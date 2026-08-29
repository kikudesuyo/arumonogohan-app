package gcp

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	v1 "github.com/kikudesuyo/arumonogohan-app/api/handler/v1"
)

func init() {
	functions.HTTP("LinebotCallback", v1.HandleLinebotCallback)
}
