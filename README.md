# Module debug-arm-viz

This module helps debug issues when the vizualization of a real arm does not match what is seen in real life. The module requires two arms.

1. a source arm which should be the real life arm we are debugging
2. a destination arm which should be a fake arm using the kinemtaics of the real arm

## Model mattmacf:debug-arm-viz:debug-arm-viz

The module has one model for a generic service that will sync the joint positions from src -> dst arm every second. This helps detect out of sync issues between real life hardware and our idealistic fakes. The model also has a DoCommand `log` which will spit out useful info on `JointPositions` and `Geometries` for both arms at that point in time.

### Dependencies

This model requires two arm resources to work (the src and target)

### Configuration

The following attribute template can be used to configure this model:

```json
{
"src_arm_name": <string>,
"dst_arm_name": <string>
}
```

#### Example Configuration

```json
{
  "src_arm_name": "arm",
  "dst_arm_name": "fake-arm"
}
```

### DoCommand

The model accepts a single `log` DoCommand which outputs helpful information to aid in the debugging process of improper arm visualization

#### Example DoCommand

```json
{
  "log": true
}
```
