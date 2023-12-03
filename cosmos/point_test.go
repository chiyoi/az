package cosmos

import (
	"context"
	"fmt"

	"github.com/chiyoi/az/identity"
	"github.com/chiyoi/iter/res"
)

const (
	EndpointCosmos = "https://neko03cosmos.documents.azure.com:443/"
)

func Example() {
	// Point Create/Read/Delete example.
	cred, err := identity.DefaultCredential()
	client, err := res.R(cred, err, NewClient(EndpointCosmos, nil))
	if err != nil {
		panic(err)
	}
	c, err := client.NewContainer("neko0001", "point_test")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	var neko, neko1 string

	exist, err := PointExist(ctx, "neko")(c)
	if err != nil {
		panic(err)
	}

	if exist {
		err = PointRead(ctx, "neko", &neko1)(c)
		if err != nil {
			panic(err)
		}

		err = PointUpdate(ctx, "neko", "nyan")(c)
		if err != nil {
			panic(err)
		}

		err = PointRead(ctx, "neko", &neko)(c)
		if err != nil {
			panic(err)
		}

		err = PointDelete(ctx, "neko")(c)
		if err != nil {
			panic(err)
		}
	} else {
		err = PointCreate(ctx, "neko", "nyan")(c)
		if err != nil {
			panic(err)
		}

		err = PointRead(ctx, "neko", &neko)(c)
		if err != nil {
			panic(err)
		}

		err = PointUpdate(ctx, "neko", "nyan1")(c)
		if err != nil {
			panic(err)
		}

		err = PointRead(ctx, "neko", &neko1)(c)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(neko)
	fmt.Println(neko1)
	// Output:
	// nyan
	// nyan1
}
