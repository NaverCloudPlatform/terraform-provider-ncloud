package ncloud

func DiffByArg(old []interface{}, new []interface{}, arg string) ([]interface{}, []interface{}) {
	var added []interface{}
	var deleted []interface{}
	oldMark := make([]bool, len(old))
	newMark := make([]bool, len(new))

	for o, oldItem := range old {
		oldItemMap := oldItem.(map[string]interface{})
		for n, newItem := range new {
			newItemMap := newItem.(map[string]interface{})
			if oldItemMap[arg].(string) == newItemMap[arg].(string) {
				oldMark[o] = true
				newMark[n] = true
			}
		}
	}

	for i, noChange := range oldMark {
		if !noChange {
			deleted = append(deleted, old[i])
		}
	}
	for i, noChange := range newMark {
		if !noChange {
			added = append(added, new[i])
		}
	}

	return added, deleted
}
