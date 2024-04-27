package object

import "math"

var Module_Math = NewModuleType("Math")

func init() {
	Module_Math.SetProperty("E", Module_Math_E)
	Module_Math.SetProperty("Pi", Module_Math_Pi)
	Module_Math.SetProperty("Phi", Module_Math_Phi)
	Module_Math.SetProperty("Sqrt2", Module_Math_Sqrt2)
	Module_Math.SetProperty("SqrtE", Module_Math_SqrtE)
	Module_Math.SetProperty("SqrtPi", Module_Math_SqrtPi)
	Module_Math.SetProperty("SqrtPhi", Module_Math_SqrtPhi)
	Module_Math.SetProperty("Ln2", Module_Math_Ln2)
	Module_Math.SetProperty("Log2E", Module_Math_Log2E)
	Module_Math.SetProperty("Ln10", Module_Math_Ln10)
	Module_Math.SetProperty("Log10E", Module_Math_Log10E)

	Module_Math.AddMethod(Module_Math_Acos)
	Module_Math.AddMethod(Module_Math_Acosh)
	Module_Math.AddMethod(Module_Math_Asin)
	Module_Math.AddMethod(Module_Math_Asinh)
	Module_Math.AddMethod(Module_Math_Atan)
	Module_Math.AddMethod(Module_Math_Atan2)
	Module_Math.AddMethod(Module_Math_Atanh)
	Module_Math.AddMethod(Module_Math_Cbrt)
	Module_Math.AddMethod(Module_Math_Cos)
	Module_Math.AddMethod(Module_Math_Cosh)
	Module_Math.AddMethod(Module_Math_Dim)
	Module_Math.AddMethod(Module_Math_Erf)
	Module_Math.AddMethod(Module_Math_Erfc)
	Module_Math.AddMethod(Module_Math_Erfcinv)
	Module_Math.AddMethod(Module_Math_Erfinv)
	Module_Math.AddMethod(Module_Math_Exp)
	Module_Math.AddMethod(Module_Math_Exp2)
	Module_Math.AddMethod(Module_Math_Expm1)
	Module_Math.AddMethod(Module_Math_FMA)
	Module_Math.AddMethod(Module_Math_Gamma)
	Module_Math.AddMethod(Module_Math_Hypot)
	Module_Math.AddMethod(Module_Math_Ilogb)
	Module_Math.AddMethod(Module_Math_J0)
	Module_Math.AddMethod(Module_Math_J1)
	Module_Math.AddMethod(Module_Math_Jn)
	Module_Math.AddMethod(Module_Math_Ldexp)
	Module_Math.AddMethod(Module_Math_Log)
	Module_Math.AddMethod(Module_Math_Log10)
	Module_Math.AddMethod(Module_Math_Log1p)
	Module_Math.AddMethod(Module_Math_Log2)
	Module_Math.AddMethod(Module_Math_Logb)
	Module_Math.AddMethod(Module_Math_Mod)
	Module_Math.AddMethod(Module_Math_Nextafter)
	Module_Math.AddMethod(Module_Math_Pow)
	Module_Math.AddMethod(Module_Math_Pow10)
	Module_Math.AddMethod(Module_Math_Sin)
	Module_Math.AddMethod(Module_Math_Sinh)
	Module_Math.AddMethod(Module_Math_Sqrt)
	Module_Math.AddMethod(Module_Math_Tan)
	Module_Math.AddMethod(Module_Math_Tanh)
	Module_Math.AddMethod(Module_Math_Y0)
	Module_Math.AddMethod(Module_Math_Y1)
	Module_Math.AddMethod(Module_Math_Yn)
	Module_Math.AddMethod(Module_Math_Deg)
	Module_Math.AddMethod(Module_Math_Rad)
	Module_Math.AddMethod(Module_Math_Primes)
}

var Module_Math_E = NewNumber(2.71828182845904523536028747135266249775724709369995957496696763)       // https://oeis.org/A001113
var Module_Math_Pi = NewNumber(3.14159265358979323846264338327950288419716939937510582097494459)      // https://oeis.org/A000796
var Module_Math_Phi = NewNumber(1.61803398874989484820458683436563811772030917980576286213544862)     // https://oeis.org/A001622
var Module_Math_Sqrt2 = NewNumber(1.41421356237309504880168872420969807856967187537694807317667974)   // https://oeis.org/A002193
var Module_Math_SqrtE = NewNumber(1.64872127070012814684865078781416357165377610071014801157507931)   // https://oeis.org/A019774
var Module_Math_SqrtPi = NewNumber(1.77245385090551602729816748334114518279754945612238712821380779)  // https://oeis.org/A002161
var Module_Math_SqrtPhi = NewNumber(1.27201964951406896425242246173749149171560804184009624861664038) // https://oeis.org/A139339
var Module_Math_Ln2 = NewNumber(0.693147180559945309417232121458176568075500134360255254120680009)    // https://oeis.org/A002162
var Module_Math_Log2E = NewNumber(1 / 0.693147180559945309417232121458176568075500134360255254120680009)
var Module_Math_Ln10 = NewNumber(2.30258509299404568401799145468436420760110148862877297603332790) // https://oeis.org/A002392
var Module_Math_Log10E = NewNumber(1 / 2.30258509299404568401799145468436420760110148862877297603332790)

func makeMathFn(fn func(float64) float64, name string, desc string) *BuiltinFunction {
	return F(
		func(scope *Scope, args ...Object) Object {
			return NewNumber(fn(args[0].(*Number).Value))
		},
		name,
		desc,
		P("x", V.Type(NumberId)),
	)
}

func makeMathFn2(fn func(float64, float64) float64, name string, desc string) *BuiltinFunction {
	return F(
		func(scope *Scope, args ...Object) Object {
			return NewNumber(fn(args[0].(*Number).Value, args[1].(*Number).Value))
		},
		name,
		desc,
		P("x", V.Type(NumberId)),
		P("y", V.Type(NumberId)),
	)
}

func makeMathFn2Inv(fn func(float64, float64) float64, name string, desc string) *BuiltinFunction {
	return F(
		func(scope *Scope, args ...Object) Object {
			return NewNumber(fn(args[0].(*Number).Value, args[1].(*Number).Value))
		},
		name,
		desc,
		P("y", V.Type(NumberId)),
		P("x", V.Type(NumberId)),
	)
}

func makeMathFn3(fn func(float64, float64, float64) float64, name string, desc string) *BuiltinFunction {
	return F(
		func(scope *Scope, args ...Object) Object {
			return NewNumber(fn(args[0].(*Number).Value, args[1].(*Number).Value, args[2].(*Number).Value))
		},
		name,
		desc,
		P("x", V.Type(NumberId)),
		P("y", V.Type(NumberId)),
		P("z", V.Type(NumberId)),
	)
}

var Module_Math_Acos = makeMathFn(math.Acos, "Acos", "Acos returns the arccosine, in radians, of x.")
var Module_Math_Acosh = makeMathFn(math.Acosh, "Acosh", "Acosh returns the inverse hyperbolic cosine of x.")
var Module_Math_Asin = makeMathFn(math.Asin, "Asin", "Asin returns the arcsine, in radians, of x.")
var Module_Math_Asinh = makeMathFn(math.Asinh, "Asinh", "Asinh returns the inverse hyperbolic sine of x.")
var Module_Math_Atan = makeMathFn(math.Atan, "Atan", "Atan returns the arctangent, in radians, of x.")
var Module_Math_Atan2 = makeMathFn2Inv(math.Atan2, "Atan2", "Atan2 returns the arctangent of y/x, using the signs of the two to determine the quadrant of the result.")
var Module_Math_Atanh = makeMathFn(math.Atanh, "Atanh", "Atanh returns the inverse hyperbolic tangent of x.")
var Module_Math_Cbrt = makeMathFn(math.Cbrt, "Cbrt", "Cbrt returns the cube root of x.")
var Module_Math_Cos = makeMathFn(math.Cos, "Cos", "Cos returns the cosine of the radian argument x.")
var Module_Math_Cosh = makeMathFn(math.Cosh, "Cosh", "Cosh returns the hyperbolic cosine of x.")
var Module_Math_Dim = makeMathFn2(math.Dim, "Dim", "Dim returns the maximum of x-y or 0.")
var Module_Math_Erf = makeMathFn(math.Erf, "Erf", "Erf returns the error function of x.")
var Module_Math_Erfc = makeMathFn(math.Erfc, "Erfc", "Erfc returns the complementary error function of x.")
var Module_Math_Erfcinv = makeMathFn(math.Erfcinv, "Erfcinv", "Erfcinv returns the inverse of Erfc(x).")
var Module_Math_Erfinv = makeMathFn(math.Erfinv, "Erfinv", "Erfinv returns the inverse error function of x.")
var Module_Math_Exp = makeMathFn(math.Exp, "Exp", "Exp returns e**x, the base-e exponential of x.")
var Module_Math_Exp2 = makeMathFn(math.Exp2, "Exp2", "Exp2 returns 2**x, the base-2 exponential of x.")
var Module_Math_Expm1 = makeMathFn(math.Expm1, "Expm1", "Expm1 returns e**x - 1, the base-e exponential of x minus 1.")
var Module_Math_FMA = makeMathFn3(math.FMA, "FMA", "FMA returns x * y + z, computed with only one rounding.")
var Module_Math_Gamma = makeMathFn(math.Gamma, "Gamma", "Gamma returns the Gamma function of x.")
var Module_Math_Hypot = makeMathFn2(math.Hypot, "Hypot", "Hypot returns Sqrt(x*x + y*y), taking care to avoid overflow and underflow.")
var Module_Math_Ilogb = makeMathFn(func(x float64) float64 { return float64(math.Ilogb(x)) }, "Ilogb", "Ilogb returns the binary exponent of x.")
var Module_Math_J0 = makeMathFn(math.J0, "J0", "J0 returns the order-zero Bessel function of the first kind.")
var Module_Math_J1 = makeMathFn(math.J1, "J1", "J1 returns the order-one Bessel function of the first kind.")
var Module_Math_Jn = makeMathFn2(func(x float64, y float64) float64 { return math.Jn(int(x), y) }, "Jn", "Jn returns the order-n Bessel function of the first kind.")
var Module_Math_Ldexp = makeMathFn2(func(x float64, y float64) float64 { return math.Ldexp(x, int(y)) }, "Ldexp", "Ldexp is the inverse of Frexp. It sets x to mantissa * 2**exp and returns x.")
var Module_Math_Log = makeMathFn(math.Log, "Log", "Log returns the natural logarithm of x.")
var Module_Math_Log10 = makeMathFn(math.Log10, "Log10", "Log10 returns the decimal logarithm of x.")
var Module_Math_Log1p = makeMathFn(math.Log1p, "Log1p", "Log1p returns the natural logarithm of 1 plus x.")
var Module_Math_Log2 = makeMathFn(math.Log2, "Log2", "Log2 returns the binary logarithm of x.")
var Module_Math_Logb = makeMathFn(math.Logb, "Logb", "Logb returns the binary exponent of x.")
var Module_Math_Mod = makeMathFn2(math.Mod, "Mod", "Mod returns the floating-point remainder of x/y.")
var Module_Math_Nextafter = makeMathFn2(math.Nextafter, "Nextafter", "Nextafter returns the next representable float64 value after x in the direction of y.")
var Module_Math_Pow = makeMathFn2(math.Pow, "Pow", "Pow returns x**y, the base-x exponential of y.")
var Module_Math_Pow10 = makeMathFn(func(x float64) float64 { return math.Pow10(int(x)) }, "Pow10", "Pow10 returns 10**n, the base-10 exponential of n.")
var Module_Math_Sin = makeMathFn(math.Sin, "Sin", "Sin returns the sine of the radian argument x.")
var Module_Math_Sinh = makeMathFn(math.Sinh, "Sinh", "Sinh returns the hyperbolic sine of x.")
var Module_Math_Sqrt = makeMathFn(math.Sqrt, "Sqrt", "Sqrt returns the square root of x.")
var Module_Math_Tan = makeMathFn(math.Tan, "Tan", "Tan returns the tangent of the radian argument x.")
var Module_Math_Tanh = makeMathFn(math.Tanh, "Tanh", "Tanh returns the hyperbolic tangent of x.")
var Module_Math_Y0 = makeMathFn(math.Y0, "Y0", "Y0 returns the order-zero Bessel function of the second kind.")
var Module_Math_Y1 = makeMathFn(math.Y1, "Y1", "Y1 returns the order-one Bessel function of the second kind.")
var Module_Math_Yn = makeMathFn2(func(x float64, y float64) float64 { return math.Yn(int(x), y) }, "Yn", "Yn returns the order-n Bessel function of the second kind.")

var Module_Math_Deg = makeMathFn(func(x float64) float64 { return x * 180 / math.Pi }, "Deg", "Deg converts x from radians to degrees.")
var Module_Math_Rad = makeMathFn(func(x float64) float64 { return x * math.Pi / 180 }, "Rad", "Rad converts x from degrees to radians.")

var Module_Math_Primes = F(
	func(scope *Scope, args ...Object) Object {
		D := map[int]int{}
		q := 0
		return NewInternalStream(func(s *Scope) Object {
			if q == 0 {
				q = 1
				return YieldWith(Two)
			}

			for {
				q += 2
				p, ok := D[q]
				if !ok {
					delete(D, q)
					D[q*q] = q
					return YieldWith(NewNumber(float64(q)))
				} else {
					x := q + 2*p
					for _, ok := D[x]; ok; {
						x += 2 * p
						_, ok = D[x]
					}
					D[x] = p
				}
			}
		}, scope)
	},
	"Primes",
	"Generate infinite prime numbers.",
)
