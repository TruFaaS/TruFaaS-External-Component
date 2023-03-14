package fission

type Function struct {
	FunctionInformation FunctionInformation `json:"function_information"`
	PackageInformation  PackageInformation  `json:"package_information"`
}

type (
	FunctionInformation struct {
		Name      string       `json:"function_name"`
		Namespace string       `json:"function_namespace"`
		Spec      FunctionSpec `json:"function_spec"`
	}

	PackageInformation struct {
		Name      string      `json:"package_name"`
		Namespace string      `json:"package_namespace"`
		Spec      PackageSpec `json:"package_spec"`
	}

	FunctionSpec struct {
		Environment     Environment    `json:"environment"`
		PackageRef      PackageRef     `json:"package_ref"`
		InvokeStrategy  InvokeStrategy `json:"invoke_strategy"`
		FunctionTimeout int            `json:"function_timeout"`
		IdleTimeout     int            `json:"idle_timeout"`
		Concurrency     int            `json:"concurrency"`
		RequestsPerPod  int            `json:"requests_per_pod"`
	}

	PackageRef struct {
		Namespace       string `json:"namespace"`
		Name            string `json:"name"`
		ResourceVersion string `json:"resource_version"`
	}

	InvokeStrategy struct {
		ExecutionStrategy ExecutionStrategy `json:"execution_strategy"`
		StrategyType      string            `json:"strategy_type"`
	}

	ExecutionStrategy struct {
		ExecutorType          string `json:"executor-type"`
		MinScale              int    `json:"min_scale"`
		MaxScale              int    `json:"max_scale"`
		TargetCPUPercent      int    `json:"target_cpu_percent"`
		SpecializationTimeout int    `json:"specialization_timeout"`
	}

	PackageSpec struct {
		Environment Environment `json:"environment"`
		Source      Archive     `json:"source,omitempty"`
		Deployment  Archive     `json:"deployment,omitempty"`
		Buildcmd    string      `json:"buildcmd,omitempty"`
	}

	Environment struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}

	Archive struct {
		Type     string   `json:"type,omitempty"`
		Literal  []byte   `json:"literal,omitempty"`
		URL      string   `json:"url,omitempty"`
		Checksum Checksum `json:"checksum,omitempty"`
	}

	Checksum struct {
		Type string `json:"type,omitempty"`
		Sum  string `json:"sum,omitempty"`
	}
)
