package go_cielo_conecta

import "fmt"

// POST /1/physicalSales/
func (c *Client) Authorization(s *Sale) (*Sale, error) {
	salePayed := &Sale{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.env.APIUrl, "/1/physicalSales/"), s)
	if err != nil {
		return salePayed, err
	}

	err = c.Send(req, &salePayed)
	if err != nil {
		return salePayed, err
	}

	return salePayed, nil
}
