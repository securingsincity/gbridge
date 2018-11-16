package gbridge

type DeviceType string

const (
	DeviceTypeAcUnit       DeviceType = "action.devices.types.AC_UNIT"
	DeviceTypeAirpurifier  DeviceType = "action.devices.types.AIRPURIFIER"
	DeviceTypeCamera       DeviceType = "action.devices.types.CAMERA"
	DeviceTypeCoffeeMaker  DeviceType = "action.devices.types.COFFEE_MAKER"
	DeviceTypeDishwasher   DeviceType = "action.devices.types.DISHWASHER"
	DeviceTypeDryer        DeviceType = "action.devices.types.DRYER"
	DeviceTypeFan          DeviceType = "action.devices.types.FAN"
	DeviceTypeKettle       DeviceType = "action.devices.types.KETTLE"
	DeviceTypeLight        DeviceType = "action.devices.types.LIGHT"
	DeviceTypeOutlet       DeviceType = "action.devices.types.OUTLET"
	DeviceTypeOven         DeviceType = "action.devices.types.OVEN"
	DeviceTypeRefrigerator DeviceType = "action.devices.types.REFRIGERATOR"
	DeviceTypeScene        DeviceType = "action.devices.types.SCENE"
	DeviceTypeSprinkler    DeviceType = "action.devices.types.SPRINKLER"
	DeviceTypeSwitch       DeviceType = "action.devices.types.SWITCH"
	DeviceTypeThermostat   DeviceType = "action.devices.types.THERMOSTAT"
	DeviceTypeVacuum       DeviceType = "action.devices.types.VACUUM"
	DeviceTypeWasher       DeviceType = "action.devices.types.WASHER"
)

type DeviceTrait string

const (
	DeviceTraitCameraStream       DeviceTrait = "action.devices.traits.CameraStream"
	DeviceTraitColorSetting       DeviceTrait = "action.devices.traits.ColorSetting"
	DeviceTraitDock               DeviceTrait = "action.devices.traits.Dock"
	DeviceTraitFanSpeed           DeviceTrait = "action.devices.traits.FanSpeed"
	DeviceTraitLocator            DeviceTrait = "action.devices.traits.Locator"
	DeviceTraitModes              DeviceTrait = "action.devices.traits.Modes"
	DeviceTraitOnOff              DeviceTrait = "action.devices.traits.OnOff"
	DeviceTraitRunCycle           DeviceTrait = "action.devices.traits.RunCycle"
	DeviceTraitScene              DeviceTrait = "action.devices.traits.Scene"
	DeviceTraitStartStop          DeviceTrait = "action.devices.traits.StartStop"
	DeviceTraitTemperatureControl DeviceTrait = "action.devices.traits.TemperatureControl"
	DeviceTraitTemperatureSetting DeviceTrait = "action.devices.traits.TemperatureSetting"
	DeviceTraitToggles            DeviceTrait = "action.devices.traits.Toggles"
)

type DeviceError string

const (
	DeviceErrorAuthExpired     DeviceError = "authExpired"
	DeviceErrorAuthFailure     DeviceError = "authFailure"
	DeviceErrorDeviceOffline   DeviceError = "deviceOffline"
	DeviceErrorTimeout         DeviceError = "timeout"
	DeviceErrorDeviceTurnedOff DeviceError = "deviceTurnedOff"
	DeviceErrorDeviceNotFound  DeviceError = "deviceNotFound"
	DeviceErrorValueOutofRange DeviceError = "valueOutOfRange"
	DeviceErrorNotSupported    DeviceError = "notSupported"
	DeviceErrorProtocolError   DeviceError = "protocolError"
	DeviceErrorUnknownError    DeviceError = "unknownError"
)

type DeviceName struct {
	DefaultNames []string `json:"defaultNames"`
	Name         string   `json:"name"`
	Nicknames    []string `json:"nicknames"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HwVersion    string `json:"hwVersion"`
	SwVersion    string `json:"swVersion"`
}

type Device struct {
	Id              string        `json:"id"`
	Type            DeviceType    `json:"type"`
	Traits          []DeviceTrait `json:"traits"`
	Name            DeviceName    `json:"name"`
	WillReportState bool          `json:"willReportState"`
	Attributes      struct {
	} `json:"attributes,omitempty"`
	RoomHint   string                 `json:"roomHint,omitempty"`
	DeviceInfo *DeviceInfo            `json:"deviceInfo,omitempty"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

func (b *Bridge) HandleExec(d Device, execFn ExecHandlerFunc) {
	if b.Devices == nil {
		b.Devices = make(map[string]*DeviceContext)
	}
	if _, ok := b.Devices[d.Id]; !ok {
		b.Devices[d.Id] = new(DeviceContext)
	}
	b.Devices[d.Id].Device = d
	b.Devices[d.Id].Exec = execFn
}

func (b *Bridge) HandleQuery(d Device, queryFn QueryHandlerFunc) {
	if b.Devices == nil {
		b.Devices = make(map[string]*DeviceContext)
	}
	if _, ok := b.Devices[d.Id]; !ok {
		b.Devices[d.Id] = new(DeviceContext)
	}
	b.Devices[d.Id].Device = d
	b.Devices[d.Id].Query = queryFn
}
