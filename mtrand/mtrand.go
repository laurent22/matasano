package mtrand

// Mersenne twister implementation based on pseudo-code at:
// http://en.wikipedia.org/wiki/Mersenne_twister#Pseudocode

type Generator struct {
	mt []int
	index int
}

func NewGenerator() *Generator {
	output := new(Generator)
	output.index = 0
	return output
}

func (this *Generator) Initialize(seed int) {
	this.index = 0
	this.mt = make([]int, 624)
	this.mt[0] = seed
	for i := 1; i < len(this.mt); i++ {
		t := (1812433253 * (this.mt[i-1] ^ (this.mt[i-1] >> 30)) + i); 
		this.mt[i] = t & 0xffffffff;
	}
}

func (this *Generator) GetInt() int {
	if this.index == 0 {
		this.generateNumbers()
	}

	y := this.mt[this.index]
	y = y ^ (y >> 11)
	y = y ^ ((y << 7) & 0x9d2c5680)
	y = y ^ ((y << 15) & 0xefc60000)
	y = y ^ ((y >> 18))

	this.index = (this.index + 1) % 624
	
	return y
 }

func (this *Generator) generateNumbers() {
	for i := 0; i < 624; i++ {
		y := (this.mt[i] & 0x80000000) + (this.mt[(i+1) % 624] & 0x7fffffff) 
		this.mt[i] = this.mt[(i + 397) % 624] ^ (y >> 1)
		if y % 2 != 0 {
			this.mt[i] = this.mt[i] ^ 2567483615
		}
	}
}