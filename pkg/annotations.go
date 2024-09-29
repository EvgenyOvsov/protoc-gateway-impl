package pkg

import (
	"google.golang.org/genproto/googleapis/api/annotations"
)

func FindGoogleAnnotationURL(ext *annotations.HttpRule) string {
	switch v := ext.Pattern.(type) {
	case *annotations.HttpRule_Get:
		return v.Get
	case *annotations.HttpRule_Post:
		return v.Post
	case *annotations.HttpRule_Put:
		return v.Put
	case *annotations.HttpRule_Delete:
		return v.Delete
	case *annotations.HttpRule_Patch:
		return v.Patch
	default:
		return ""
	}
}
