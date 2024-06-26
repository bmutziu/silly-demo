// Code generated by timoni. DO NOT EDIT.
// Note that this file is required and should contain
// the values schema and the timoni workflow.

package main

import (
	templates "bmutziu.me/silly-demo/templates"
)

values: templates.#Config

timoni: {
	apiVersion: "v1alpha1"
	instance: templates.#Instance & {
		config: values
		config: metadata: {
			name:      string @tag(name)
			namespace: string @tag(namespace)
		}
	}
	apply: all: [ for obj in instance.objects {obj}]
}
