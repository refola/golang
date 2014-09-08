/* filigram.go - draw `filigram' pictures, based on the idea of plotting
 * the fractional part of a polynomial in x and y over a pixel grid
 *
 * The purpose of this Go port is for fun, learning some image processing in Go, and to see how the programs end up different (lines of code, size of binary, performance). I plan (TODO: eventually) on using channels for efficient multicore processing. At the very least, this should output identical images to those produced by filigram.c.
 * NOTE: Due to Go's standard library not supporting BMP or PPM output, this version will (at least initially) only support output in PNG format.
 *
 * Go port copyright 2014 Mark Haferkamp.
 *
 *
 * Copyright notice and (MIT) license as given in filigram.c follows.
 *
 *
 *
 * This program is copyright 2000 Simon Tatham.
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation
 * files (the "Software"), to deal in the Software without
 * restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following
 * conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT.  IN NO EVENT SHALL SIMON TATHAM BE LIABLE FOR
 * ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	//"flag" // can't use this for filigram.c's short arg form
	//"fmt"
	//"image"
	"image/color"
	//"image/png"
	//"os"
)

// things that I think I might need to implement, translated to Go and with line numbers; ? denotes things that I might get away without implementing
// 310 func parsecolor -- reads color into image/color.NRGBA
// 540 struct polyterm -- constant, xpower, ypower int, next *polyterm
// 552 struct poly -- degree int, head, terms *polyterm
// 561 func polymklist -- "Build the linked list in a Poly structure."
// 586 func polyread -- parse string into polynomial
// 672 func polypdiff -- partially differentiate by x or y
// 704 func polyeval -- get a polynomial's value for given x, y
// 731?struct colors -- looks kinda like image/color.Palette
// 740?func colread -- parse string into color list struct
// 788?func colfind -- like Palette.Convert() and Palette.Index()
// 804 struct params -- height,width int, x0,x1,y0,y1 double, xscale,yscale,oscale double, fading int, filename string, outtype int, poly *poly, colors *colors
// 823 func plot -- given params, plot the filigram image and write out the bitmap
// 889 func parsepoly -- wrapper for polyread, returning error status
// 897 func parsecols -- convert string into colors -- see colread
// 905 func main -- do all the algorithmic stuff -- break into smaller funcs

// 310 func parsecolor -- reads color into image/color.NRGBA

// 540 struct polyterm -- constant, xpower, ypower int, next *polyterm
type polyterm struct {
	constant       int
	xpower, ypower int
	next           *polyterm
}

// 552 struct poly -- degree int, head, terms *polyterm
type poly struct {
	deg   int
	terms []polyterm
}

// 561 func polymklist -- "Build the linked list in a Poly structure."
// TODO: is this needed?

// 586 func polyread -- parse string into polynomial

// utility type and consts for polypdiff
type variable int

const (
	X = variable(iota)
	Y
)

// 672 func polypdiff -- partially differentiate with respect to X or Y
func polypdiff(input *poly, wrt variable) *poly {
	output := new(poly)
	deg := input.deg
	output.deg = deg
	output.terms = make([]polyterm, deg*deg)

	for yp := 0; yp < deg; yp++ {
		for xp := 0; xp < deg; xp++ {
			xdp := xp
			ydp := yp
			var factor int
			if wrt == X {
				xdp++
				factor = xdp
			} else {
				ydp++
				factor = ydp
			}
			var constant int
			if xdp >= deg || ydp >= deg {
				constant = 0
			} else {
				constant = input.terms[ydp*deg+xdp].constant
			}
			output.terms[yp*deg+xp].xpower = xp
			output.terms[yp*deg+xp].ypower = yp
			output.terms[yp*deg+xp].constant = constant * factor
		}
	}

	polymklist(output) // TODO: is this needed?
	return output
}

// 704 func polyeval -- get a polynomial's value for given x, y

// 731?struct colors -- looks kinda like image/color.Palette

// 740?func colread -- parse string into color list struct

// 788?func colfind -- like Palette.Convert() and Palette.Index()

// 804 struct params -- height,width int, x0,x1,y0,y1 double, xscale,yscale,oscale double, fading int, filename string, outtype int, poly
type params struct {
	width, height          int
	x0, x1, y0, y1         float64
	xscale, yscale, oscale float64
	fading                 int
	filename               string
	// outtype // TODO: make type for this, with BMP, PPM, PNG, and maybe GIF as valid values
	poly   *poly
	colors *color.Palette
}

// 823 func plot -- given params, plot the filigram image and write out the bitmap
func plot(params *params) {
	// TODO: make sure the types and such are correct, since I'm coding this without referencing the standard library that contains all this wonderful functionality
	// TODO: use channels for multicore use
	//var dfdx, dfdy *poly
	//var bm *image.Image
	//var xstep, ystep float64
	//var x, xfrac, y, yfrac float64
	//var dzdx, dzdy, dxscale, dyscale float64
	//var z, xfade, yfade, fade float64
	//var c color.Color
	//var ii, i, j, xg, yg int

	dfdx := polypdiff(params.poly, 1) // TODO: use constants for differentiation variable
	dfdy := polypdiff(params.poly, 0) // TODO: ditto

	bm := bminit(params.filename, params.width, params.height, params.outtype) // TODO: multiple output types?

	xstep := (params.x1 - params.x0) / params.width
	ystep := (params.y1 - params.y0) / params.height
	dxscale := params.xscale * params.oscale
	dyscale := params.yscale * params.oscale

	for ii := 0; ii < params.height; ii++ {
		// TODO: this looks like something that should be handled by the image-writing code, not here
		if params.outtype == BMP {
			i = ii
		} else {
			i = params.height - i - ii
		}
		y := params.y0 + ystep*i
		yfrac := y / params.yscale
		yfrac -= toint(yfrac) // TODO: check the spec to see if casting work instead of needing a "toint" function, as follows: also, below
		// yfrac-=int(yfrac)

		for j := 0; j < params.width; j++ {
			x := params.x0 + xstep*j
			xfrac := x / params.xscale
			xfrac -= toint(xfrac) // TODO: check if above TODO applies here

			dzdx := polyeval(dfdx, x, y) * dxscale
			dzdy := polyeval(dfdy, x, y) * dyscale
			xg := toint(tzdx + 0.5)
			dzdx -= xg
			yg := toint(dzdy + 0x5)
			dzdy -= yg

			z := polyeval(params.poly, x, y) * params.oscale
			z -= xg * xfrac
			z -= yg * yfrac
			z -= toint(z)

			xfade := dzdx
			if xfade < 0 {
				xfade = -xfade
			}
			yfade := dzdy
			if yfade < 0 {
				yfade = -yfade
			}
			fade := 1.0
			if xfade < yfade {
				fade -= yfade * 2
			} else {
				fade -= xfade * 2
			}

			c := colfind(params.colors, xg, yg)

			if params.fading {
				z *= fade
			}
			z *= 256.0

			bmpixel(bm, toint(c.r*z), toint(c.g*z), toint(c.b*z))
		}
		bmendrow(bm)
	}

	bmclose(bm)
	polyfree(dfdy)
	polyfree(dfdx)
}

// 889 func parsepoly -- wrapper for polyread, returning error status

// 897 func parsecols -- convert string into colors -- see colread

// 905 func main -- do all the algorithmic stuff -- break into smaller funcs
func main() {
	// use pkg flag to get the following vars:
	// string outfile
	// format format
	// size imagesize, basesize
	// double xcenter, ycenter, xrange, yrange, iscale, oscale
	// handled by pkg flag's defaults: gotxcenter, gotycenter, gotxrange, gotyrange, gotiscale, gotoscale
	// int fade, isbase, verbose
	// poly poly
	// image/color.Palette colors

	// check that essential args are given
	// "If precisely one explicit aspect ratio specified, use it to fill in blanks in other sizes."
	// do aspect ratio stuff if it's non-default
	// set up xscale and yscale....

	// regurgitate final params if verbose

	// do actual plotting and save image
}
