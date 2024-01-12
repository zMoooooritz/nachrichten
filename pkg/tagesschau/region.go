package tagesschau

type RegionID int

const (
	DE RegionID = iota
	BW
	BY
	BE
	BB
	HB
	HH
	HE
	MV
	NI
	NW
	RP
	SL
	SN
	ST
	SH
	TH
)

type RegionName string

var GERMAN_NAMES = map[RegionID]RegionName{
	DE: "Deutschland",
	BW: "Baden-Württemberg",
	BY: "Bayern",
	BE: "Berlin",
	BB: "Brandenburg",
	HB: "Bremen",
	HH: "Hamburg",
	HE: "Hessen",
	MV: "Mecklenburg-Vorpommern",
	NI: "Niedersachsen",
	NW: "Nordrhein-Westfalen",
	RP: "Rheinland-Pfalz",
	SL: "Saarland",
	SN: "Sachsen",
	ST: "Sachsen-Anhalt",
	SH: "Schleswig-Holstein",
	TH: "Thüringen",
}

var ENGLISH_NAMES = map[RegionID]RegionName{
	DE: "Germany",
	BW: "Baden-Württemberg",
	BY: "Bavaria",
	BE: "Berlin",
	BB: "Brandenburg",
	HB: "Bremen",
	HH: "Hamburg",
	HE: "Hessen",
	MV: "Mecklenburg-Western Pomerania",
	NI: "Lower Saxony",
	NW: "North Rhine-Westphalia",
	RP: "Rhineland-Palatinate",
	SL: "Saarland",
	SN: "Saxony",
	ST: "Saxony-Anhalt",
	SH: "Schleswig-Holstein",
	TH: "Thuringia",
}
