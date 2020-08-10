# tfplanparse

tfplanparse is a Go library for parsing `terraform plan` outputs. Supports `terraform` v0.12 CLI output currently. Does not fully support multiline diffs right now.

**NOTE:** This library does _not_ parse the file produced by `terraform plan -out=<file>`. If you want to parse that file, you should use [`hashicorp/terraform-json`](https://github.com/hashicorp/terraform-json).

This project is still WIP and the schemas are subject to change.

## Why

By parsing the output from `terraform plan` into something more machine readable, you can perform some validation and workflows depending on what the output contains. Some example use cases include:

- Filter out unexpected changes from routine `terraform` operations
- Send a notification when a certain type of resource is being modified
- Trigger additional workflows before `terraform apply` if a certain resource type is modified

## `tfplanparse` vs `terraform-json`

You can parse the output from `terraform plan` using `terraform-json` by running `terraform plan -out=<file>` followed by `terraform state show -json <file>`. However, while the output file is encoded, it contains sensitive information, and `terraform state show -json <file>` will [output this sensitive information in plain text](https://www.terraform.io/docs/commands/show.html).

With `tfplanparse`, you can parse the output from `terraform plan` directly, skipping writing a state file containing sensitive data to disk. Additionally, sensitive data is redacted in the `terraform plan` output, so the secrets will remain redacted when parsed.

## Usage

```go
import "github.com/drlau/tfplanparse"

func main() {
    result, err := tfplanparse.Parse(os.Stdin)
    if err != nil {
        panic(err)
    }

    // or, you can read directly from a file
    fromFile, err := tfplanparse.ParseFromFile("out.stdout")
    if err != nil {
        panic(err)
    }
}
```

The returned type from `Parse` and `ParseFromFile` is `[]*tfplanparse.ResourceChange`. Each `ResourceChange` corresponds to a single resource in the `terraform plan` output and has the following fields:

- **`Address`**: Absolute resource address
- **`ModuleAddress`**: Module portion of the absolute address, if any
- **`Type`**: The type of the resource (example: `gcp_instance.foo` -> `"gcp_instance"`)
- **`Name`**: The name of the resource (example: `gcp_instance.foo` -> `"foo"`)
- **`Index`**: The index key for resources created with `count` or `for_each`
- **`UpdateType`**: The type of update (refer to `updatetype.go` for possible values)
- **`Tainted`**: Indicates whether the resource is tainted or not
- **`AttributeChanges`**: Planned attribute changes
- **`MapAttributeChanges`**: Planned attribute changes that are maps

Each `ResourceChange` also has the following helper functions:

- **`GetBeforeResource`**: Returns a copy of the resource before the planned changes as a `map[string]interface{}`
- **`GetAfterResource`**: Returns a copy of the resource after the planned changes as a `map[string]interface{}`

Additionally, these helper functions accept the following options:

- **`IgnoreComputed`**
- **`IgnoreSensitive`**
- **`ComputedOnly`**