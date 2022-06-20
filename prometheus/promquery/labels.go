package promquery

type L map[string]string

func (l L) Clone() L {
	newL := make(L, len(l))
	for k, v := range l {
		newL[k] = v
	}
	return newL
}

func (l L) Merge(l2 L) L {
	newL := l.Clone()
	for k, v := range l2 {
		newL[k] = v
	}
	return newL
}

func LabelsKV(labelsKV ...string) L {
	labels := make(L, len(labelsKV)/2)
	for i := 0; i < len(labelsKV)/2; i++ {
		labels[labelsKV[i]] = labelsKV[i+1]
	}
	return labels
}
