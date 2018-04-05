# fproto-wrap-validator-std

[fproto-wrap-validator](https://github.com/RangelReale/fproto-wrap-validator) standard validator.

### example


```protobuf
syntax = "proto3";
package gw_sample;
option go_package = "gwsample/core";

import "github.com/RangelReale/fproto-wrap-validator-std/validator.proto";

message User {
    // name is required (must not be blank)
    string name = 1 [(validator.field) = {required: true}];
    // the length of the REPEATED array must be greater then zero (validator.rfield), and each of the array values
    // cannot be blank and their length must be greater then 10.
    repeated string email = 2 [(validator.field) = {required: true, length_gt: 10}, (validator.rfield) = {length_gt: 0}];
}
```

### field validators

* (validator.field).required [bool]: field must not be empty. The concept of "empty" varies depending of the field type.
* (validator.field).regex [string]: field string value must match the regex.
* (validator.field).string_eq [string]: field string value must be *exactly* this value.

* (validator.field).int_gt [int64]: field integer value must be greater than this value.
* (validator.field).int_lt [int64]: field integer value must be lower than this value.
* (validator.field).int_gte [int64]: field integer value must be greater or equals this value.
* (validator.field).int_lte [int64]: field integer value must be lower or equals this value.
* (validator.field).int_eq [int64]: field integer value must be *exactly* this value.
* (validator.field).int_enum_check [bool]: if field type is enum, checks if the integer value is among the declared ones.

* (validator.field).float_epsilon [double]: float value tolerance (epsilon) for all "float_" checks.
* (validator.field).float_gt [double]: field float value must be greater than this value.
* (validator.field).float_lt [double]: field float value must be lower than this value.
* (validator.field).float_gte [double]: field float value must be greater or equals this value.
* (validator.field).float_lte [double]: field float value must be lower or equals this value.
* (validator.field).float_eq [double]: field float value must be *exactly* (within epsilon) this value.

* (validator.field).length_gt [int64]: field length must be greater than this value.
* (validator.field).length_lt [int64]: field length must be lower than this value.
* (validator.field).length_gte [int64]: field length must be greater or equals this value.
* (validator.field).length_lte [int64]: field length must be lower or equals this value.
* (validator.field).length_eq [int64]: field length must be *exactly* this value.

* (validator.field).bool_eq [bool]: bool field must be *exactly* this value.

### repeated field validators

* (validator.rfield).required [bool]: item list must have at least one item.

* (validator.rfield).length_gt [int64]: item list length must be greater than this value.
* (validator.rfield).length_lt [int64]: item list length must be lower than this value.
* (validator.rfield).length_gte [int64]: item list length must be greater or equals this value.
* (validator.rfield).length_lte [int64]: item list length must be lower or equals this value.
* (validator.rfield).length_eq [int64]: item list length must be *exactly* this value.

### author

Rangel Reale (rangelspam@gmail.com)

