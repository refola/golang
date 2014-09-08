/* Plan
Implement a program to convert between TI-format and plain text/XML files using the downloaded file specs.

File support priorities:
Ext	Type
8Xn	real number
8Xs	string
8Xp	program
*/

/* TI-83+ file format
Offset	Length	Name		Description
0	8	signature	always "**TI83F*".
8	3	further sig	always contains {1Ah, 0Ah, 00h} = {26, 10, 0}
11	42	Comment		either zero-terminated or right padded with spaces.
53	2	Date length	should be 57 bytes less than the file size.
55	n	Data section	a number of variable entries
55+n	2	Checksum	lower 16 bits of sum of all bytes in data section.	Note - All 2-byte integers are stored little-endian Intel-style (least significant byte first).
*/

/* Data format
Offset	Length	Name		Description
0	2	Header type	Bh if Version and Flag skipped or Dh if present.
2	2	Length		number of bytes of variable data.
4	1	type ID		see Type IDs
5	8	Variable name	padded with NULL characters (0h) on the right.
13	1	Version		usually set to 0 (present if first bytes are Dh).
14	1	Flag		80h if archived, 00h else (present if first bytes are Dh).
15	2	Length		number of bytes for variable (copy of value in offset 2)
17	n	Data		see Variable formats
*/

/* Type IDs and file extensions
IDs from page 53 of http://education.ti.com/downloads/guidebooks/sdk/83p/sdk83pguide.pdf	Value	Type
00h	Real
01h	List
02h	Matrix
03h	Equation
04h	String
05h	Program
06h	Protected Program
07h	Picture
08h	Graph Database
0Bh	New EQU Obj
0Ch	Complex Obj
0Dh	Complex List Obj
14h	Application Obj
15h	AppVar Obj
17h	Group Obj
ID	Ext	Type

Normal stuff
00h	8Xn	real number
0Ch	8Xc	complex number
01h	8Xl	real list
0Dh	8Xl	complex list
02h	8Xm	matrix
03h	8Xy	Y-Variable (equation)
04h	8Xs	string
08h	8Xd	GDB (function, polar, parametric or sequence)
07h	8Xi	picture (image)
05h	8Xp	program
06h	8Xp	protected program
	8Xw	window settings (Window or RclWindow)
	8Xz	zoom (saved window settings)
	8Xt	table setup
Flash
	8Xk	FLASH application
	8Xq	FLASH certificate
	8Xu	FLASH Operating System
	8Xv	Application Variable
Groups
	8Xg	Multiple variables of varying types (group)
	8Xgrp	TiLP only: 'group' variable
*/

// Variable formats

/* Real Numbers
Offset	Length	Name		Description
0	1	Flags		(see table below)
1	1	Exponent	in base-10
2	7	Mantissa	14-digit unsigned binary-coded-decimal number	Bit	Description
1	If set, number is undefined (used for initial sequence values)
2	If bits 2 and 3 are set and bit 1 is clear, the number is half of a complex variable.
3
6	Uncertain. Likely if set, number has not been modified since last graph.
7	Sign bit: set if number is negative.
*/

/* Complex numbers
Offset	Length	Description
0	9	A real number describing the "real" component of the complex number.
9	9	A real number describing the "imaginary" component of the complex number.
*/

/* Lists
Offset	Length	Description
0	2	Number of elements in the list
1	n	Element values, first to last. Each element is a 9-byte real number.
*/

/* Matrices
Offset	Length	Description
0	1	Number of columns in the matrix (no more than 255)
1	1	Number of rows in the matrix (no more than 255)
1	n	Element values, one by one. Each element is a 9-byte real number.	The element values are arranged in row definitions from top to bottom. Each row consists of a number of real or complex elements from left to right.
*/

/* Y-Variables
Offset	Length	Description
0	2	Number of token bytes in the Y-Variable. Note that some tokens use two bytes.
2	n	Tokens, first to last.
*/

/* Strings
Offset	Length	Description
0	2	Number of token bytes. Note that some tokens use two bytes.
2	n	Tokens, first to last.
*/

/* Graphics Databases
function-mode
Offset	Length	Description
0	2	Length, in bytes, of GDB, minus two.
2	1	Unknown
3	1	Graphing Mode ID. This byte has a value of 10h for function GDB's.
4	1	Mode settings (see mode setting table below)
5	1	Unused - has a value of 80h.
6	1	Extended mode settings (see extended mode setting table below)
7	9	A real number: Xmin
16	9	A real number: Xmax
25	9	A real number: Xscl
34	9	A real number: Ymin
43	9	A real number: Ymax
52	9	A real number: Yscl
61	9	A real number: Xres
70	10	Ten style bytes, for Y1-Y9 and Y0, respectively (see style table below).
80	n	Ten functions, for Y1-Y9 and Y0, respectively (see function table below).

parametric-mode
Offset	Length	Description
0	2	Length, in bytes, of GDB, minus two.
2	1	Unknown - has a value of 0h.
3	1	Graphing Mode ID. This byte has a value of 40h for parametric GDB's.
4	1	Mode settings (see mode setting table below)
5	1	Unused - has a value of 80h.
6	1	Extended mode settings (see extended mode setting table below)
7	9	A real number: Xmin
16	9	A real number: Xmax
25	9	A real number: Xscl
34	9	A real number: Ymin
43	9	A real number: Ymax
52	9	A real number: Yscl
61	9	A real number: Tmin
70	9	A real number: Tmax
79	9	A real number: Tstep
70	6	Six style bytes, for X1T/Y1T-X6T/Y6T, respectively (see style table below).
76	n	Twelve functions, for X1T-X6T and Y1T-Y6T, respectively (see function table).

polar-mode
Offset	Length	Description
0	2	Length, in bytes, of GDB, minus two.
2	1	Unknown - has a value of 0h.
3	1	Graphing Mode ID. This byte has a value of 20h for polar GDB's.
4	1	Mode settings (see mode setting table below)
5	1	Unused - has a value of 80h.
6	1	Extended mode settings (see extended mode setting table below)
7	9	A real number: Xmin
16	9	A real number: Xmax
25	9	A real number: Xscl
34	9	A real number: Ymin
43	9	A real number: Ymax
52	9	A real number: Yscl
61	9	A real number: ϴmin
70	9	A real number: ϴmax
79	9	A real number: ϴstep
70	6	Six style bytes, for r1-r6, respectively (see style table below).
76	n bytes	Six functions, for r1-r6, respectively (see function table below).

sequence-mode
Offset	Length	Description
0	2	Length, in bytes, of GDB, minus two.
2	1	Unknown - has a value of 0h.
3	1	Graphing Mode ID. This byte has a value of 80h for sequence GDB's.
4	1	Mode settings (see mode setting table below)
5	1	Sequence mode settings (see sequence mode setting table below)
6	1	Extended mode settings (see extended mode setting table below)
7	9	A real number: Xmin
16	9	A real number: Xmax
25	9	A real number: Xscl
34	9	A real number: Ymin
43	9	A real number: Ymax
52	9	A real number: Yscl
61	9	A real number: PlotStart
70	9	A real number: nMax
79	9	A real number: u(nMin), first element
88	9	A real number: v(nMin), first element
97	9	A real number: nMin
106	9	A real number: u(nMin), second element
115	9	A real number: v(nMin), second element
124	9	A real number: w(nMin), first element
133	9	A real number: PlotStep
142	9	A real number: w(nMin), second element
151	3	Three style bytes, for u, v and w, respectively (see style table below).
154	n	Three functions, for u, v, and w, respectively (see function table below).

mode setting byte following format
Bit (Mask)	Mode if set (1)	Mode if clear (0)
0	Dot	Connected
1	Simul	Sequential
2	GridOn	GridOff
3	PolarGC	RectGC
4	CoordOff	CoordOn
5	AxesOff	AxesOn
6	LabelOn	LabelOff
7	This bit is always clear.	extended mode setting byte
Bit (Mask)	Mode if set (1)	Mode if clear (0)
0	ExprOff	ExprOn	sequence mode setting byte
Bit	Mode if set	Mode if clear
0	Web		Time, uv, vw or uw
1	This bit is always clear.
2	uv		Time, web, vw or uw
3	vw		Time, web, uv or uw
4	uw		Time, web, uv or vw
5	These bits are always clear.
6	These bits are always clear.
7	This bit is always set.	style byte format
Value	Graph Style
0	[solid line]
1	[thick line]
2	[shade above]
3	[shade below]
4	[trace]
5	[animate]
6	[dotted line]	function definition format
Offset	Length	Description
0	1	Flags - 23h if function is selected or 03h if deselected or undefined
1	n	A Y-Variable defining the function. n=0 for undefined functions.
*/

/* Pictures
Offset	Length	Description
0	2	Size of picture data (always 2F4h)
2	1008	1-bit-per-pixel bitmap, left-to-right and top-to-bottom, rows before columns
*/

/* Programs
Offset	Length	Description
0	2	Number of token bytes in the program. Note that some tokens use two bytes.
2	n	Tokens, first to last.	Asm programs
0	2	Length
2	n	ASCII-encoded hexadecimal digits, ending with tokenized ":End\n:0000\n:End"
Is edit-locked or edit-unlocked, depending on the type ID.
*/

/* Window Settings
normal
Offset	Length	Description
0	2	Always has a value of D0h.
2	1	Unknown - value is 00h.
3	9	A real number: Xmin
12	9	A real number: Xmax
21	9	A real number: Xscl
30	9	A real number: Ymin
39	9	A real number: Ymax
48	9	A real number: Yscl
57	9	A real number: ϴmin
66	9	A real number: ϴmax
75	9	A real number: ϴstep
84	9	A real number: Tmin
93	9	A real number: Tmax
102	9	A real number: Tstep
111	9	A real number: PlotStart
120	9	A real number: nMax
129	9	A real number: u(nMin), first element
138	9	A real number: v(nMin), first element
147	9	A real number: nMin
156	9	A real number: u(nMin), second element
165	9	A real number: v(nMin), second element
174	9	A real number: w(nMin), first element
183	9	A real number: PlotStep
192	9	A real number: Xres
201	9	A real number: w(nMin), second element

saved
Offset	Length	Description
0	2	Always has a value of CFh.
2	9	A real number: Xmin
11	9	A real number: Xmax
20	9	A real number: Xscl
29	9	A real number: Ymin
38	9	A real number: Ymax
47	9	A real number: Yscl
56	9	A real number: ϴmin
65	9	A real number: ϴmax
74	9	A real number: ϴstep
83	9	A real number: Tmin
92	9	A real number: Tmax
101	9	A real number: Tstep
110	9	A real number: PlotStart
119	9	A real number: nMax
128	9	A real number: u(nMin), first element
137	9	A real number: v(nMin), first element
146	9	A real number: nMin
155	9	A real number: u(nMin), second element
164	9	A real number: v(nMin), second element
173	9	A real number: w(nMin), first element
182	9	A real number: PlotStep
191	9	A real number: Xres
200	9	A real number: w(nMin), second element
*/

/* Table Settings
Offset	Length	Description
0	2	Always has a value of 12h.
2	9	A real number: TblMin
10	9	A real number: ΔTbl
*/

// Token meanings

/* First bytes
Row, Col: First digit, second digit.
Underscore characters signify space characters.

 	0	1	2	3	4	5	6	7	8	9	A	B	C	D	E	F
0	███████	▸DMS	▸Dec	▸Frac	→	BoxPlot	[	]	{	}	r	°	-1	2	T	3	1
(	)	round(	pxl-Test(	augment(	rowSwap(	row+(	*row(	*row+(	max(	min(	R▸Pr(	R▸Pϴ(	P▸Rx(	P▸Ry(	median(	2
randM(	mean(	solve(	seq(	fnInt(	nDeriv(	 	fMin(	fMax(	_	"	,	i	!	CubicReg_	QuartReg_	3
0	1	2	3	4	5	6	7	8	9	.	 E 	_or_	_xor_	:	note	4
and	A	B	C	D	E	F	G	H	I	J	K	L	M	N	O	5
P	Q	R	S	T	U	V	W	X	Y	Z		more	more	more	prgm	6
more	more	more	more	Radian	Degree	Normal	Sci	Eng	Float	=	<	>				7
+	-	Ans	Fix_	Horiz	Full	Func	Param	Polar	Seq	IndpntAuto	IndpntAsk	DependAuto	DependAsk	more		8
*	/	Trace	ClrDraw	ZStandard	ZTrig	ZBox	Zoom_In	Zoom_Out	ZSquare	ZInteger	ZPrevious	ZDecimal	ZoomStat	9
ZoomRcl	PrintScreen	ZoomSto	Text(	_nPr_	_nCr_	FnOn_	FnOff_	StorePic_	RecallPic_	StoreGDB_	RecallGDB_	Line(	Vertical_	Pt-On(	Pt-Off(	A
Pt-Change(	Pxl-On(	Pxl-Off(	Pxl-Change(	Shade(	Circle(	Horizontal_	Tangent(	DrawInv_	DrawF_	more	rand		getKey	'	?	B
-	int(	abs(	det(	identity(	dim(	sum(	prod(	not(	iPart(	fPart(	more	(	3(	ln_(	e^(	C
log(	10^(	sin(	sin-1(	cos(	cos-1(	tan(	tan-1(	sinh(	sinh-1(	cosh(	cosh-1(	tanh(	tanh-1(	If_	Then	D
Else	While_	Repeat_	For(	End	Return	Lbl_	Goto_	Pause_	Stop	IS>(	DS>(	Input_	Prompt_	Disp_	DispGraph	E
Output(	ClrHome	Fill(	SortA(	SortD(	DispTable	Menu(	Send(	Get(	PlotsOn_	PlotsOff_	L	Plot1(	Plot2(	Plot3(	 	F
^		1-Var_Stats_	2-Var_Stats_	LinReg(a+bx)_	ExpReg_	LnReg_	PwrReg_	Med-Med_	QuadReg_	ClrList_	ClrTable	Histogram	xyLine	Scatter	LinReg(ax+b)_	*/
/* System variables
 	5C	5D	5E	60	61	62	63	AA
00
[A]	L1	 	Pic1	GDB1	 	ZXscl	Str1	01
[B]	L2	 	Pic2	GDB2	RegEq	ZYscl	Str2	02
[C]	L3	 	Pic3	GDB3	n	Xscl 	Str3	03
[D]	L4	 	Pic4	GDB4		Yscl	Str4	04
[E]	L5	 	Pic5	GDB5	x	UnStart	Str5	05
[F]	L6	 	Pic6	GDB6	x2	VnStart	Str6	06
[G]	L7	 	Pic7	GDB7	Sx	Un-1	Str7	07
[H]	L8	 	Pic8	GDB8	x	Vn-1	Str8	08
[I]	L9	 	Pic9	GDB9	minX	ZUnStart	Str9	09
[J]	L0	 	Pic0	GDB0	maxX	ZVnStart	Str0	0A
 	 	 	 	 	minY	Xmin	 	0B
 	 	 	 	 	maxY	Xmax	 	0C
 	 	 	 	 		Ymin	 	0D
 	 	 	 	 	y	Ymax	 	0E
 	 	 	 	 	y2	Tmin	 	0F
 	 	 	 	 	Sy	Tmax	 	10
 	 	Y1	 	 	y	min	 	11
 	 	Y2	 	 	xy	max	 	12
 	 	Y3	 	 	r	ZXmin	 	13
 	 	Y4	 	 	Med	ZXmax	 	14
 	 	Y5	 	 	Q1	ZYmin	 	15
 	 	Y6	 	 	Q3	ZYmax	 	16
 	 	Y7	 	 	a	Zmin	 	17
 	 	Y8	 	 	b	Zmax	 	18
 	 	Y9	 	 	c	ZTmin	 	19
 	 	Y0	 	 	d	ZTmax	 	1A
 	 	 	 	 	e	TblMin	 	1B
 	 	 	 	 	x1	nMin	 	1C
 	 	 	 	 	x2	ZnMin	 	1D
 	 	 	 	 	x3	nMax	 	1E
 	 	 	 	 	y1	ZnMax	 	1F
 	 	 	 	 	y2	nStart	 	20
 	 	X1T	 	 	y3	ZnStart	 	21
 	 	Y1T	 	 	n	Tbl	 	22
 	 	X2T	 	 	p	Tstep	 	23
 	 	Y2T	 	 	z	step	 	24
 	 	X3T	 	 	t	ZTstep	 	25
 	 	Y3T	 	 	2	Zstep	 	26
 	 	X4T	 	 		X	 	27
 	 	Y4T	 	 	df	Y	 	28
 	 	X5T	 	 		XFact	 	29
 	 	Y5T	 	 	1	YFact	 	2A
 	 	X6T	 	 	2	TblInput 	2B
 	 	Y6T	 	 	1	N	 	2C
 	 	 	 	 	Sx1	I%	 	2D
 	 	 	 	 	n1	PV	 	2E
 	 	 	 	 	2	PMT	 	2F
 	 	 	 	 	Sx2	FV	 	20
 	 	 	 	 	n2	Xres	 	21
 	 	 	 	 	Sxp	ZXres	 	22
 	 	 	 	 	lower	 	 	23
 	 	 	 	 	upper	 	 	24
 	 	 	 	 	s	 	 	25
 	 	 	 	 	r2	 	 	26
 	 	 	 	 	R2	 	 	27
 	 	 	 	 	df	 	 	28
 	 	 	 	 	SS	 	 	29
 	 	 	 	 	MS	 	 	2A
 	 	 	 	 	df	 	 	2B
 	 	 	 	 	SS	 	 	2C
 	 	 	 	 	MS	 	 	40
 	 	r1	 	 	 	 	 	41
 	 	r2	 	 	 	 	 	42
 	 	r3	 	 	 	 	 	43
 	 	r4	 	 	 	 	 	44
 	 	r5	 	 	 	 	 	45
 	 	r6	 	 	 	 	 	80
 	 	u	 	 	 	 	 	81
 	 	v
*/
package documentation
