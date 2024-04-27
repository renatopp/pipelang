data Vector {
  x = 0
  y = 0
  
  fn String(this) {
    return sprintf('<%d, %d>', this.x, this.y)
  }
}

vec := Vector { y = 100 }
vec.String()
