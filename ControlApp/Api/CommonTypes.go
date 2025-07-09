package Api

type Color struct {
	R  int
	G  int
	B  int
	W  int
	A  int
	UV int
}

func isColorValid(color Color) bool {
	return color.R >= 0 && color.R <= 0xFF &&
		color.G >= 0 && color.G <= 0xFF &&
		color.B >= 0 && color.B <= 0xFF &&
		color.W >= 0 && color.W <= 0xFF &&
		color.A >= 0 && color.A <= 0xFF &&
		color.UV >= 0 && color.UV <= 0xFF
}
