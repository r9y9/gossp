/*************************************************************************/
/*                                                                       */
/*                  Language Technologies Institute                      */
/*                     Carnegie Mellon University                        */
/*                         Copyright (c) 2013                            */
/*                        All Rights Reserved.                           */
/*                                                                       */
/*  Permission is hereby granted, free of charge, to use and distribute  */
/*  this software and its documentation without restriction, including   */
/*  without limitation the rights to use, copy, modify, merge, publish,  */
/*  distribute, sublicense, and/or sell copies of this work, and to      */
/*  permit persons to whom this work is furnished to do so, subject to   */
/*  the following conditions:                                            */
/*   1. The code must retain the above copyright notice, this list of    */
/*      conditions and the following disclaimer.                         */
/*   2. Any modifications must be clearly marked as such.                */
/*   3. Original authors' names are not deleted.                         */
/*   4. The authors' names are not used to endorse or promote products   */
/*      derived from this software without specific prior written        */
/*      permission.                                                      */
/*                                                                       */
/*  CARNEGIE MELLON UNIVERSITY AND THE CONTRIBUTORS TO THIS WORK         */
/*  DISCLAIM ALL WARRANTIES WITH REGARD TO THIS SOFTWARE, INCLUDING      */
/*  ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO EVENT   */
/*  SHALL CARNEGIE MELLON UNIVERSITY NOR THE CONTRIBUTORS BE LIABLE      */
/*  FOR ANY SPECIAL, INDIRECT OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES    */
/*  WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN   */
/*  AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION,          */
/*  ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF       */
/*  THIS SOFTWARE.                                                       */
/*                                                                       */
/*************************************************************************/
/*             Author:  Alok Parlikar (aup@cs.cmu.edu)                   */
/*               Date:  January 2013                                     */
/*                                                                       */
/*  The code slightly modified by Ryuichi Yamamoto, Feb.2 2014           */
/*************************************************************************/

package mcepalpha

import (
	"math"
)

// Calculate all-pass constant (alpha) for a given sample frequency.
func CalcMcepAlpha(sampfreq int) float64 {
	var alphastart, alphastep, alphaend float64

	// Number of points to sample from the mel scale and warped frequency curves
	numpoints := 1000

	alphastart = 0
	alphastep = 0.001
	alphaend = 1

	minDistance := math.MaxFloat64
	bestAlpha := 0.0

	mel := MelScaleVector(sampfreq, numpoints)

	for alpha := alphastart; alpha < alphaend; alpha += alphastep {
		distance := RMSDistance(mel, WarpingVector(alpha, numpoints))
		if distance < minDistance {
			minDistance = distance
			bestAlpha = alpha
		}
	}
	//fmt.Printf("Best Alpha for Frequency %d: %f\n", sampfreq, bestAlpha)

	return bestAlpha
}

// Generate a Mel Scale for sampling frequency _sampfreq_ and return a
// normalized vector of length _vectorlength_ containing equally
// spaced points between 0 and (sampfreq/2)
func MelScaleVector(sampfreq int, vectorlength int) []float64 {
	var i int
	var step float64
	var melscalevector []float64

	step = (float64(sampfreq) / 2.0) / float64(vectorlength)

	melscalevector = make([]float64, vectorlength, vectorlength)

	for i = 0; i < vectorlength; i++ {
		var melscale float64

		// Equations taken from Wikipedia
		f := step * float64(i)
		melscale = (1000.0 / math.Log(2)) * math.Log(1.0+(f/1000.0))

		melscalevector[i] = melscale
	}

	// Normalize the Vector. Values are already positive, and
	// monotonically increasing, so divide by max
	max := melscalevector[vectorlength-1]
	for i = 0; i < vectorlength; i++ {
		melscalevector[i] /= max
	}

	return melscalevector
}

// Generate a warped frequency curve for warping coefficient _alpha_
// and return a normalized vector of length _vectorlength_ containing
// equally spaced points between 0 and PI
func WarpingVector(alpha float64, vectorlength int) []float64 {
	var i int
	var step float64
	var warpingvector []float64

	step = math.Pi / float64(vectorlength)

	warpingvector = make([]float64, vectorlength, vectorlength)

	for i = 0; i < vectorlength; i++ {
		var warpfreq, omega float64
		var num, den float64
		omega = step * float64(i)
		num = (1 - alpha*alpha) * math.Sin(omega)
		den = (1+alpha*alpha)*math.Cos(omega) - 2*alpha
		warpfreq = math.Atan(num / den)

		// Wrap the phase if it becomes negative
		if warpfreq < 0 {
			warpfreq += math.Pi
		}

		warpingvector[i] = warpfreq
	}

	// Normalize the Vector. Values are positive and monotonically
	// increasing, so divide by max
	max := warpingvector[vectorlength-1]
	for i = 0; i < vectorlength; i++ {
		warpingvector[i] /= max
	}
	return warpingvector
}

// Get RootMeanSquare distance between two vectors
// This assumes that vectors are normalized
func RMSDistance(vector1, vector2 []float64) float64 {
	var sum, count float64

	sum = 0
	count = 0

	for i, val1 := range vector1 {
		val2 := vector2[i]
		sum += (val2 - val1) * (val2 - val1)
		count += 1
	}

	return math.Sqrt(sum / count)
}
