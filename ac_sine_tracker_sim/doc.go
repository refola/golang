/*
AC sine wave simulator and tracker for algorithmic testing

Requirements

Simulate a noisy 60 Hz sine wave and use the output to track the zeros. Then call an interrupt routine to tell the device whenever zeros happen.

Wave constraints

	* Must approximate 50 or 60 Hz (configurable Hz?)
	* Wave may shift up to 3% in a cycle, but probably not that much
	* Things can be noisy, maybe ±10% of max amplitude?
	* Must output 600,000 times a second (10,000 times per wave)

Tracker constraints

	* Must use very little memory and processing
	* May give bad results during particularly noisy conditions
	* Must eventually correct itself
	* Has to run an interrupt when and only when the wave passes through the zero point

*/
package main
