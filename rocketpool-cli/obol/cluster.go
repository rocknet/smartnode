package obol

import (
	"fmt"

	"github.com/urfave/cli"
)

func getCharonHealth(c *cli.Context) error {
	// // Get RP client
	// rp, err := rocketpool.NewClientFromCtx(c).WithReady()
	// if err != nil {
	// 	return err
	// }
	// defer rp.Close()

	// response, err := rp.GetCharonHealth()
	// if err != nil {
	// 	return fmt.Errorf("Error fetching charon health: %w", err)
	// }
	// // Log & return
	// fmt.Println("Successfully fetched charon health: ", response)
	fmt.Println("Successfully fetched charon health: TODO")
	return nil
}
