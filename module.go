package debugarmviz

import (
	"context"
	"fmt"
	"time"

	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	DebugArmViz = resource.NewModel("mattmacf", "debug-arm-viz", "debug-arm-viz")
)

func init() {
	resource.RegisterService(generic.API, DebugArmViz,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newDebugArmVizDebugArmViz,
		},
	)
}

type Config struct {
	SrcArmName string `json:"src_arm_name"`
	DstArmName string `json:"dst_arm_name"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	if cfg.SrcArmName == "" {
		return nil, nil, fmt.Errorf("src_arm_name is required")
	}
	if cfg.DstArmName == "" {
		return nil, nil, fmt.Errorf("dst_arm_name is required")
	}
	return []string{cfg.SrcArmName, cfg.DstArmName}, nil, nil
}

type debugArmVizDebugArmViz struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	srcArm arm.Arm
	dstArm arm.Arm
}

func newDebugArmVizDebugArmViz(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewDebugArmViz(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func (s *debugArmVizDebugArmViz) syncArms(ctx context.Context) error {
	fps := 1
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(time.Duration(1000/fps) * time.Millisecond)
		}

		srcPositions, err := s.srcArm.JointPositions(ctx, map[string]interface{}{})
		if err != nil {
			return err
		}
		err = s.dstArm.MoveToJointPositions(ctx, srcPositions, map[string]interface{}{})
		if err != nil {
			return err
		}

	}
}

func NewDebugArmViz(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	srcArm, err := arm.FromDependencies(deps, conf.SrcArmName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	dstArm, err := arm.FromDependencies(deps, conf.DstArmName)
	if err != nil {
		cancelFunc()
		return nil, err
	}
	s := &debugArmVizDebugArmViz{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		srcArm:     srcArm,
		dstArm:     dstArm,
	}

	go s.syncArms(cancelCtx)
	return s, nil
}

func (s *debugArmVizDebugArmViz) Name() resource.Name {
	return s.name
}

func (s *debugArmVizDebugArmViz) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	_, ok := cmd["log"]
	if ok {
		err := s.logArmInfo(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"success": true,
		}, nil
	}
	return nil, fmt.Errorf("unknown command")
}

func (s debugArmVizDebugArmViz) logArmInfo(ctx context.Context) error {
	srcPositions, err := s.srcArm.JointPositions(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	srcGeometries, err := s.srcArm.Geometries(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	dstPositions, err := s.dstArm.JointPositions(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	dstGeometries, err := s.dstArm.Geometries(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	s.logger.Infow("srcPositions: %v", srcPositions)
	s.logger.Infow("srcGeometries: %v", srcGeometries)
	s.logger.Infow("dstPositions: %v", dstPositions)
	s.logger.Infow("dstGeometries: %v", dstGeometries)
	return nil
}

func (s *debugArmVizDebugArmViz) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
