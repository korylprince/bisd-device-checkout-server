package api

import (
	"context"
	"math"
	"strconv"
)

// ChargeURLBase is the base URL used for charge links
var ChargeURLBase = "/charges/edit?type=id&search="

// DeviceURLBase is the base URL used for device links
var DeviceURLBase = "/edit?type=id&search="

// LinkType is the type of a link
type LinkType string

// Link types
const (
	LinkTypeDevice LinkType = "device"
	LinkTypeCharge LinkType = "charge"
)

// Issue represents an issue with a student
type Issue struct {
	Description    string   `json:"description,omitempty"`
	Link           string   `json:"link,omitempty"`
	LinkType       LinkType `json:"link_type,omitempty"`
	LinkValue      float32  `json:"link_value,omitempty"`
	LinkAdditional string   `json:"link_additional,omitempty"`
}

// StatusType is the type of Chromebook a student will receive
type StatusType string

// Status types
const (
	StatusTypeNone     StatusType = "none"
	StatusTypeRedBag   StatusType = "red_bag"
	StatusTypeBlackBag StatusType = "black_bag"
)

// Status represents the status of a student
type Status struct {
	Type   StatusType `json:"type"`
	Issues []*Issue   `json:"issues,omitempty"`
}

// Status returns the Status of the student
func (s *Student) Status(ctx context.Context) (*Status, error) {
	status := &Status{Issues: make([]*Issue, 0)}

	if s.T2E2Status == nil {
		status.Type = StatusTypeNone
		status.Issues = append(status.Issues, &Issue{Description: "T2E2 Agreement not completed"})
	}

	//check for devices checked out
	devices, err := getDeviceList(ctx, s.Name())
	if err != nil {
		return nil, err
	}

	if len(devices) > 0 {
		status.Type = StatusTypeNone
		for _, d := range devices {
			status.Issues = append(status.Issues, &Issue{
				Description: "Student has device checked out",
				Link:        DeviceURLBase + strconv.Itoa(d),
				LinkType:    LinkTypeDevice,
			})
		}
	}

	//check for charges
	charges, err := getChargeList(ctx, s.Name())
	if err != nil {
		return nil, err
	}

	var noneCharges []*Charge
	var redCharges []*Charge

	for _, c := range charges {
		//subtract one to take rounding errors out of the mix
		if math.Abs(float64(c.AmountCharged()-c.AmountPaid)) < 1 {
			// charge is paid
			continue
		} else if c.AmountPaid >= (c.AmountCharged()/2)-1 {
			// half or more is paid
			redCharges = append(redCharges, c)
		} else if c.AmountCharged()-c.AmountPaid-1 > 0 {
			// less than half is paid
			noneCharges = append(noneCharges, c)
		}
	}

	if len(noneCharges) == 0 && len(redCharges) == 0 {
		if s.T2E2Status != nil && *(s.T2E2Status) == "No" {
			if status.Type != StatusTypeNone {
				status.Type = StatusTypeRedBag
			}
			status.Issues = append(status.Issues, &Issue{
				Description: "T2E2 Agreement does not permit student to take home device",
			})
		} else {
			if status.Type != StatusTypeNone {
				status.Type = StatusTypeBlackBag
			}
		}
		return status, nil
	}

	if len(noneCharges) == 0 && status.Type != StatusTypeNone {
		status.Type = StatusTypeRedBag
	} else {
		status.Type = StatusTypeNone
	}

	for _, c := range noneCharges {
		status.Issues = append(status.Issues, &Issue{
			Description:    "Student has charge with less than 50% paid",
			Link:           ChargeURLBase + strconv.Itoa(c.ID),
			LinkType:       LinkTypeCharge,
			LinkValue:      c.AmountCharged() - c.AmountPaid,
			LinkAdditional: c.Description(),
		})

	}

	for _, c := range redCharges {
		status.Issues = append(status.Issues, &Issue{
			Description:    "Student has unpaid charge",
			Link:           ChargeURLBase + strconv.Itoa(c.ID),
			LinkType:       LinkTypeCharge,
			LinkValue:      c.AmountCharged() - c.AmountPaid,
			LinkAdditional: c.Description(),
		})
	}

	return status, nil
}
