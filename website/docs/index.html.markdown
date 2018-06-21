---
layout: "ncloud"
page_title: "Provider: Local"
sidebar_current: "docs-ncloud-index"
description: |-
  The Local provider is used to manage ncloud resources, such as files.
---

# Ncloud Provider

The Ncloud provider is used to manage ncloud resources, such as files.

Use the navigation to the left to read about the available resources.

~> **Note** Terraform primarily deals with remote resources which are able
to outlive a single Terraform run, and so ncloud resources can sometimes violate
its assumptions. The resources here are best used with care, since depending
on ncloud state can make it hard to apply the same Terraform configuration on
many different ncloud systems where the ncloud resources may not be universally
available. See specific notes in each resource for more information.
