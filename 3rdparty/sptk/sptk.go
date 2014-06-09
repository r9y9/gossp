// Package sptk provides Go interface to SPTK (Speech Signal Processing Toolkit),
// which is originally written in C.
package sptk

// #cgo pkg-config: SPTK
// #include <stdio.h>
// #include <SPTK/SPTK.h>
import "C"

func ACep(audioSample float64, m int, lambda, step, tau float64, pd int,
	eps float64) ([]float64, float64) {
	resultBuffer := make([]float64, m)
	var predictionError C.double
	predictionError = C.acep(
		C.double(audioSample),
		(*C.double)(&resultBuffer[0]),
		C.int(m),
		C.double(lambda),
		C.double(step),
		C.double(tau),
		C.int(pd),
		C.double(eps))
	return resultBuffer, float64(predictionError)
}

func MFCC(audioBuffer []float64, sampleRate int, alpha, eps float64,
	wlng, flng, m, n, ceplift int, dftmode, usehamming bool) []float64 {
	// Convert go bool to C.Boolean
	var dftmodeInGo, usehammingInGo C.Boolean
	if dftmode {
		dftmodeInGo = 1
	} else {
		dftmodeInGo = 0
	}
	if usehamming {
		usehammingInGo = 1
	} else {
		usehammingInGo = 0
	}

	resultBuffer := make([]float64, m)

	C.mfcc(
		(*C.double)(&audioBuffer[0]),
		(*C.double)(&resultBuffer[0]),
		C.double(sampleRate),
		C.double(alpha),
		C.double(eps),
		C.int(wlng),
		C.int(flng),
		C.int(m),
		C.int(n),
		C.int(ceplift),
		dftmodeInGo,
		usehammingInGo)

	return resultBuffer
}

func UELS(audioBuffer []float64, flng, m, itr1, itr2 int, dd float64,
	etype int, e float64, itype int) []float64 {
	resultBuffer := make([]float64, m+1)

	C.uels(
		(*C.double)(&audioBuffer[0]),
		C.int(flng),
		(*C.double)(&resultBuffer[0]),
		C.int(m),
		C.int(itr1),
		C.int(itr2),
		C.double(dd),
		C.int(etype),
		C.double(e),
		C.int(itype))

	return resultBuffer
}

func UELSWithDefaultParameters(audioBuffer []float64, order int) []float64 {
	return UELS(audioBuffer, len(audioBuffer),
		order, 2, 30, 0.001, 0, 0.0, 0)
}

func MCep(audioBuffer []float64, flng, m int, a float64, itr1, itr2 int,
	dd float64, etype int, e, f float64, itype int) []float64 {
	resultBuffer := make([]float64, m+1)

	C.mcep(
		(*C.double)(&audioBuffer[0]),
		C.int(flng),
		(*C.double)(&resultBuffer[0]),
		C.int(m),
		C.double(a),
		C.int(itr1),
		C.int(itr2),
		C.double(dd),
		C.int(etype),
		C.double(e),
		C.double(f),
		C.int(itype))

	return resultBuffer
}

func MCepWithDefaultParameters(audioBuffer []float64,
	order int, alpha float64) []float64 {
	return MCep(audioBuffer, len(audioBuffer),
		order, alpha, 2, 30, 0.001, 0, 0.0, 0.00001, 0)
}

func MGCep(audioBuffer []float64, flng, m int, a, g float64, n, itr1, itr2 int,
	dd float64, etype int, e, f float64, itype, otype int) []float64 {
	resultBuffer := make([]float64, m+1)

	C.mgcep(
		(*C.double)(&audioBuffer[0]),
		C.int(flng),
		(*C.double)(&resultBuffer[0]),
		C.int(m),
		C.double(a),
		C.double(g),
		C.int(n),
		C.int(itr1),
		C.int(itr2),
		C.double(dd),
		C.int(etype),
		C.double(e),
		C.double(f),
		C.int(itype))

	if otype == 0 || otype == 1 || otype == 2 || otype == 4 {
		C.ignorm((*C.double)(&resultBuffer[0]),
			(*C.double)(&resultBuffer[0]),
			C.int(m),
			C.double(g))
	}

	if otype == 0 || otype == 2 || otype == 4 {
		C.b2mc((*C.double)(&resultBuffer[0]),
			(*C.double)(&resultBuffer[0]),
			C.int(m),
			C.double(a))
	}

	if otype == 2 || otype == 4 {
		C.gnorm((*C.double)(&resultBuffer[0]),
			(*C.double)(&resultBuffer[0]),
			C.int(m),
			C.double(g))
	}

	if otype == 4 || otype == 5 {
		for i := int(m); i >= 1; i-- {
			resultBuffer[i] *= float64(g)
		}
	}

	return resultBuffer
}

func MGCepWithDefaultParameters(audioBuffer []float64,
	order int, alpha, gamma float64) []float64 {
	otype := 0 // cepstrum
	return MGCep(audioBuffer, len(audioBuffer), order,
		alpha, gamma, len(audioBuffer)-1, 2, 30, 0.001, 0, 0.0, 0.00001, 0,
		otype)
}

func MLSADF(x float64, b []float64, m int, a float64, pd int, d []float64) float64 {
	// d must be declared outside in this function as follows:
	// d := make([]float64, 3*(pd+1)+pd*(m+2))
	// see mlsadf.c in SPTK for details.

	filterd := C.mlsadf(
		C.double(x),
		(*C.double)(&b[0]),
		C.int(m),
		C.double(a),
		C.int(pd),
		(*C.double)(&d[0]))

	return float64(filterd)
}

func MGLSADF(x float64, b []float64, m int, a float64, n int, d []float64) float64 {
	// d must be declared outside this function as follows:
	// d := make([]float64, (m+1)*n)
	// see mglsadf.c in SPTK for details.

	filterd := C.mglsadf(
		C.double(x),
		(*C.double)(&b[0]),
		C.int(m),
		C.double(a),
		C.int(n),
		(*C.double)(&d[0]))

	return float64(filterd)
}

func SWIPE(audioBuffer []float64, sampleRate int, frameShift int,
	min, max, st float64, otype int) []float64 {
	resultBuffer := make([]float64, len(audioBuffer)/frameShift+1)

	C.swipe((*C.double)(&audioBuffer[0]),
		(*C.double)(&resultBuffer[0]),
		C.int(len(audioBuffer)),
		C.int(sampleRate),
		C.int(frameShift),
		C.double(min),
		C.double(max),
		C.double(st),
		C.int(otype))

	return resultBuffer
}

func SWIPEWithDefaultParameters(audioBuffer []float64,
	sampleRate, frameShift int,
	min, max float64) []float64 {
	return SWIPE(audioBuffer, sampleRate, frameShift, min, max, 0.3, 1)
}

func GNorm(ceps []float64, gamma float64) []float64 {
	normalizedCeps := make([]float64, len(ceps))

	m := len(ceps) - 1
	C.gnorm(
		(*C.double)(&ceps[0]),
		(*C.double)(&normalizedCeps[0]),
		C.int(m),
		C.double(gamma))

	return normalizedCeps
}

func IGNorm(normalizedCeps []float64, gamma float64) []float64 {
	ceps := make([]float64, len(normalizedCeps))

	m := len(normalizedCeps) - 1
	C.ignorm(
		(*C.double)(&normalizedCeps[0]),
		(*C.double)(&ceps[0]),
		C.int(m),
		C.double(gamma))

	return ceps
}

func GC2GC(c1 []float64, gamma1 float64, m2 int, gamma2 float64) []float64 {
	c2 := make([]float64, m2+1)

	m1 := len(c1) - 1

	C.gc2gc(
		(*C.double)(&c1[0]),
		C.int(m1),
		C.double(gamma1),
		(*C.double)(&c2[0]),
		C.int(m2),
		C.double(gamma2))

	return c2
}

func MGC2MGC(c1 []float64, alpha1, gamma1 float64,
	m2 int, alpha2, gamma2 float64) []float64 {
	c2 := make([]float64, m2+1)

	m1 := len(c1) - 1

	C.mgc2mgc(
		(*C.double)(&c1[0]),
		C.int(m1),
		C.double(alpha1),
		C.double(gamma1),
		(*C.double)(&c2[0]),
		C.int(m2),
		C.double(alpha2),
		C.double(gamma2))

	return c2
}

func FreqT(ceps []float64, order int, alpha float64) []float64 {
	warpedCeps := make([]float64, order+1)

	m1 := len(ceps) - 1

	C.freqt(
		(*C.double)(&ceps[0]),
		C.int(m1),
		(*C.double)(&warpedCeps[0]),
		C.int(order),
		C.double(alpha))

	return warpedCeps
}

func MC2B(mcep []float64, alpha float64) []float64 {
	filterCoef := make([]float64, len(mcep))

	m := len(mcep) - 1

	C.mc2b(
		(*C.double)(&mcep[0]),
		(*C.double)(&filterCoef[0]),
		C.int(m),
		C.double(alpha))

	return filterCoef
}

func B2MC(b []float64, alpha float64) []float64 {
	mcep := make([]float64, len(b))

	m := len(b) - 1

	C.b2mc(
		(*C.double)(&b[0]),
		(*C.double)(&mcep[0]),
		C.int(m),
		C.double(alpha))

	return mcep
}

func C2IR(ceps []float64, length int) []float64 {
	h := make([]float64, length)

	m := len(ceps) - 1

	C.c2ir(
		(*C.double)(&ceps[0]),
		C.int(m+1),
		(*C.double)(&h[0]),
		C.int(length))

	return h
}
