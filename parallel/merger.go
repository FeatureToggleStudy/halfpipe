package parallel

import "github.com/springernature/halfpipe/manifest"

type Merger struct {
}

func NewParallelMerger() Merger {
	return Merger{}
}

func (Merger) MergeParallelTasks(tasks manifest.TaskList) (mergedTasks manifest.TaskList) {
	tmpParallel := manifest.Parallel{}
	tmpParallelName := ""
	for _, task := range tasks {
		if task.GetParallelGroup().IsSet() {
			currentParallelName := string(task.GetParallelGroup())
			if tmpParallelName != currentParallelName {
				if len(tmpParallel.Tasks) > 0 {
					mergedTasks = append(mergedTasks, tmpParallel)
				}
				tmpParallel = manifest.Parallel{}
				tmpParallelName = currentParallelName
			}

			tmpParallel.Tasks = append(tmpParallel.Tasks, task)
		} else {
			if len(tmpParallel.Tasks) > 0 {
				mergedTasks = append(mergedTasks, tmpParallel)
				tmpParallel = manifest.Parallel{}
				tmpParallelName = ""
			}
			mergedTasks = append(mergedTasks, task)
		}
	}

	if len(tmpParallel.Tasks) > 0 {
		mergedTasks = append(mergedTasks, tmpParallel)
	}

	return
}
