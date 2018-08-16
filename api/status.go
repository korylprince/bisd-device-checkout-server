package api

import (
	"context"
	"strconv"
)

//ChargeURLBase is the base URL used for charge links
var ChargeURLBase = "/charges/edit?type=id&search="

//DeviceURLBase is the base URL used for device links
var DeviceURLBase = "/edit?type=id&search="

//Status represents the type of Chromebook a student will receive
type Status struct {
	Type     string   `json:"type"`
	Reason   string   `json:"reason,omitempty"`
	LinkType string   `json:"link_type,omitempty"`
	Links    []string `json:"links,omitempty"`
}

//Status returns the Status of the student
func (s *Student) Status(ctx context.Context) (*Status, error) {
	if s.T2E2Status == nil {
		return &Status{Type: "none", Reason: "T2E2 Agreement not completed"}, nil
	}

	//check for devices checked out
	devices, err := getDeviceList(ctx, s.Name())
	if err != nil {
		return nil, err
	}

	if len(devices) > 0 {
		var links []string
		for _, d := range devices {
			links = append(links, DeviceURLBase+strconv.Itoa(d))
		}

		return &Status{
			Type:     "none",
			Reason:   "Student has device(s) checked out",
			LinkType: "device",
			Links:    links,
		}, nil
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
		if c.AmountCharged()-c.AmountPaid-1 <= 0 {
			//pass
		} else if c.AmountPaid >= (c.AmountCharged()/2)-1 {
			redCharges = append(redCharges, c)
		} else {
			noneCharges = append(noneCharges, c)
		}
	}

	if len(noneCharges) == 0 && len(redCharges) == 0 {
		if *(s.T2E2Status) == "No" {
			return &Status{Type: "red_bag", Reason: "T2E2 Agreement does not permit student to take home device"}, nil
		}
		return &Status{Type: "black_bag"}, nil
	}

	var links []string
	for _, c := range noneCharges {
		links = append(links, ChargeURLBase+strconv.Itoa(c.ID))
	}
	for _, c := range redCharges {
		links = append(links, ChargeURLBase+strconv.Itoa(c.ID))
	}

	if len(noneCharges) > 0 {
		return &Status{
			Type:     "none",
			Reason:   "Student has charge(s) with less than 50% paid",
			LinkType: "charge",
			Links:    links,
		}, nil
	}

	return &Status{
		Type:     "red_bag",
		Reason:   "Student has unpaid charge(s)",
		LinkType: "charge",
		Links:    links,
	}, nil
}
