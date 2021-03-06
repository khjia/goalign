package dna

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type TN93Model struct {
	// Parameters (for eigen values/vectors computation)
	// See https://en.wikipedia.org/wiki/Models_of_DNA_evolution#F81_model_(Felsenstein_1981)
	qmatrix    *mat.Dense
	leigenvect *mat.Dense
	val        []float64
	reigenvect *mat.Dense
}

func NewTN93Model() *TN93Model {
	return &TN93Model{
		nil,
		nil,
		nil,
		nil,
	}
}

func (m *TN93Model) InitModel(kappa1, kappa2, piA, piC, piG, piT float64) (err error) {
	m.qmatrix = mat.NewDense(4, 4, []float64{
		-(piC + kappa1*piG + piT), piC, kappa1 * piG, piT,
		piA, -(piA + piG + kappa2*piT), piG, kappa2 * piT,
		kappa1 * piA, piC, -(kappa1*piA + piC + piT), piT,
		piA, kappa2 * piC, piG, -(piA + kappa2*piC + piG),
	})
	// Normalization of Q
	norm := -piA*m.qmatrix.At(0, 0) -
		piC*m.qmatrix.At(1, 1) -
		piG*m.qmatrix.At(2, 2) -
		piT*m.qmatrix.At(3, 3)
	m.qmatrix.Apply(func(i, j int, v float64) float64 { return v / norm }, m.qmatrix)
	err = m.computeEigens()
	return
}

func (m *TN93Model) computeEigens() (err error) {
	var u mat.CDense

	// Compute eigen values, left and right eigenvectors of Q
	eigen := &mat.Eigen{}
	if ok := eigen.Factorize(m.qmatrix, mat.EigenRight); !ok {
		err = fmt.Errorf("Problem during matrix decomposition")
		return
	}

	val := make([]float64, 4)
	for i, b := range eigen.Values(nil) {
		val[i] = real(b)
	}
	eigen.VectorsTo(&u)
	reigenvect := mat.NewDense(4, 4, nil)
	leigenvect := mat.NewDense(4, 4, nil)
	reigenvect.Apply(func(i, j int, val float64) float64 { return real(u.At(i, j)) }, reigenvect)
	leigenvect.Inverse(reigenvect)

	m.leigenvect = leigenvect
	m.reigenvect = reigenvect
	m.val = val
	return
}

func (m *TN93Model) Eigens() (val []float64, leftvectors, rightvectors *mat.Dense, err error) {
	leftvectors = m.leigenvect
	rightvectors = m.reigenvect
	val = m.val
	return
}

func (m *TN93Model) Pij(i, j int, l float64) float64 {
	return -1.0
}

func (m *TN93Model) Analytical() bool {
	return false
}

func (m *TN93Model) NState() int {
	return 4
}
