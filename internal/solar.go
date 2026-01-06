package internal

import (
	"math"
	"time"
)

const (
	// Standard zenith angle for sunrise/sunset.
	// 90° + 50' (50 arc minutes = 50/60 degrees = 0.8333°)
	// This accounts for atmospheric refraction and the sun's radius.
	sunriseZenith = 90.8333
)

// CalculateTimes returns sunrise and sunset times for a given location and date.
func CalculateTimes(lat, lon float64, t time.Time) (sunrise, sunset time.Time) {
	date := t

	// Calculate Julian Day
	jd := julianDay(date)

	// Iterative calculation for more accuracy
	// First pass: rough estimate
	sunriseMinutes := timeOfTransit(jd, lat, lon, sunriseZenith, true)
	sunsetMinutes := timeOfTransit(jd, lat, lon, sunriseZenith, false)

	// Second pass: refined calculation using the rough estimate
	sunriseJD := jd + sunriseMinutes/1440.0
	sunsetJD := jd + sunsetMinutes/1440.0

	sunriseMinutes = timeOfTransit(sunriseJD, lat, lon, sunriseZenith, true)
	sunsetMinutes = timeOfTransit(sunsetJD, lat, lon, sunriseZenith, false)

	// Convert minutes since midnight UTC to time
	sunrise = minutesToTime(date, sunriseMinutes)
	sunset = minutesToTime(date, sunsetMinutes)

	return sunrise, sunset
}

// timeOfTransit calculates the time of sun transit for a given zenith angle.
// Returns minutes since midnight UTC.
func timeOfTransit(jd, lat, lon, zenith float64, rising bool) float64 {
	jc := julianDayToJulianCentury(jd)

	declination := sunDeclination(jc)

	hourAngle := hourAngleFromZenith(lat, declination, zenith)
	if !rising {
		hourAngle = -hourAngle
	}

	delta := -lon - hourAngle
	eqTime := equationOfTime(jc)
	offset := delta*4.0 - eqTime

	// Adjust for crossing midnight
	if offset < -720.0 {
		offset += 1440.0
	}

	return 720.0 + offset
}

// julianDay converts a date to Julian Day number.
func julianDay(t time.Time) float64 {
	utc := t.UTC()

	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()

	dayFraction := float64(utc.Hour())/24.0 +
		float64(utc.Minute())/(24.0*60.0) +
		float64(utc.Second())/(24.0*60.0*60.0)

	// Adjust for January and February
	if month <= 2 {
		year--
		month += 12
	}

	// Gregorian calendar adjustment
	a := year / 100
	b := 2 - a + (a / 4)

	jd := float64(int(365.25*float64(year+4716))) +
		float64(int(30.6001*float64(month+1))) +
		float64(day) + dayFraction + float64(b) - 1524.5

	return jd
}

// julianDayToJulianCentury converts Julian Day to Julian Century (from J2000.0).
func julianDayToJulianCentury(jd float64) float64 {
	return (jd - 2451545.0) / 36525.0
}

// geomMeanLongSun calculates the geometric mean longitude of the sun.
func geomMeanLongSun(jc float64) float64 {
	l0 := 280.46646 + jc*(36000.76983+0.0003032*jc)
	return math.Mod(l0, 360.0)
}

// geomMeanAnomalySun calculates the geometric mean anomaly of the sun.
func geomMeanAnomalySun(jc float64) float64 {
	return 357.52911 + jc*(35999.05029-0.0001537*jc)
}

// eccentricEarthOrbit calculates the eccentricity of Earth's orbit.
func eccentricEarthOrbit(jc float64) float64 {
	return 0.016708634 - jc*(0.000042037+0.0000001267*jc)
}

// sunEqOfCenter calculates the equation of center of the sun.
func sunEqOfCenter(jc float64) float64 {
	m := geomMeanAnomalySun(jc)
	mrad := math.Pi * m / 180.0

	sinm := math.Sin(mrad)
	sin2m := math.Sin(2 * mrad)
	sin3m := math.Sin(3 * mrad)

	c := sinm*(1.914602-jc*(0.004817+0.000014*jc)) +
		sin2m*(0.019993-0.000101*jc) +
		sin3m*0.000289

	return c
}

// sunTrueLong calculates the sun's true longitude.
func sunTrueLong(jc float64) float64 {
	l0 := geomMeanLongSun(jc)
	c := sunEqOfCenter(jc)
	return l0 + c
}

// sunApparentLong calculates the sun's apparent longitude.
func sunApparentLong(jc float64) float64 {
	trueLong := sunTrueLong(jc)
	omega := 125.04 - 1934.136*jc
	return trueLong - 0.00569 - 0.00478*math.Sin(math.Pi*omega/180.0)
}

// meanObliquityOfEcliptic calculates the mean obliquity of the ecliptic.
func meanObliquityOfEcliptic(jc float64) float64 {
	seconds := 21.448 - jc*(46.815+jc*(0.00059-jc*0.001813))
	return 23.0 + (26.0+(seconds/60.0))/60.0
}

// obliquityCorrection calculates the corrected obliquity of the ecliptic.
func obliquityCorrection(jc float64) float64 {
	e0 := meanObliquityOfEcliptic(jc)
	omega := 125.04 - 1934.136*jc
	return e0 + 0.00256*math.Cos(math.Pi*omega/180.0)
}

// sunDeclination calculates the sun's declination.
func sunDeclination(jc float64) float64 {
	oc := obliquityCorrection(jc)
	al := sunApparentLong(jc)

	sint := math.Sin(math.Pi*oc/180.0) * math.Sin(math.Pi*al/180.0)
	return 180.0 * math.Asin(sint) / math.Pi
}

// equationOfTime calculates the equation of time in minutes.
func equationOfTime(jc float64) float64 {
	epsilon := obliquityCorrection(jc)
	l0 := geomMeanLongSun(jc)
	e := eccentricEarthOrbit(jc)
	m := geomMeanAnomalySun(jc)

	y := math.Tan(math.Pi * epsilon / 360.0)
	y *= y

	sin2l0 := math.Sin(2 * math.Pi * l0 / 180.0)
	sinm := math.Sin(math.Pi * m / 180.0)
	cos2l0 := math.Cos(2 * math.Pi * l0 / 180.0)
	sin4l0 := math.Sin(4 * math.Pi * l0 / 180.0)
	sin2m := math.Sin(2 * math.Pi * m / 180.0)

	etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 -
		0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return 4.0 * 180.0 * etime / math.Pi
}

// hourAngleFromZenith calculates the hour angle for a given zenith.
func hourAngleFromZenith(lat, declination, zenith float64) float64 {
	latRad := math.Pi * lat / 180.0
	decRad := math.Pi * declination / 180.0
	zenithRad := math.Pi * zenith / 180.0

	h := (math.Cos(zenithRad) - math.Sin(latRad)*math.Sin(decRad)) /
		(math.Cos(latRad) * math.Cos(decRad))

	// Handle polar day/night (sun never rises/sets)
	if h > 1.0 {
		return 0.0
	}
	if h < -1.0 {
		return 180.0
	}

	return 180.0 * math.Acos(h) / math.Pi
}

// minutesToTime converts minutes since midnight UTC to local time.
func minutesToTime(date time.Time, minutes float64) time.Time {
	hours := int(minutes / 60.0)
	mins := int(minutes) % 60
	secs := int((minutes - float64(int(minutes))) * 60.0)

	utc := time.Date(date.Year(), date.Month(), date.Day(), hours, mins, secs, 0, time.UTC)

	return utc.In(date.Location())
}
