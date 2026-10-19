package config

import "sort"

const TopCommandSlots = 3

type rankedCommand struct {
	Command Command
	Count   int
}

func (p *Project) RecordCommandUse(alias string) {
	if p.CommandUsage == nil {
		p.CommandUsage = map[string]int{}
	}
	p.CommandUsage[alias]++
}

func (p *Project) RecordMetaUse(key string) {
	if key == "" {
		return
	}
	if p.MetaUsage == nil {
		p.MetaUsage = map[string]int{}
	}
	p.MetaUsage[key]++
}

func (p *Project) MetaUseCount(key string) int {
	if p.MetaUsage == nil {
		return 0
	}
	return p.MetaUsage[key]
}

func (p *Project) ranked() []rankedCommand {
	out := make([]rankedCommand, 0, len(p.Commands))
	for _, c := range p.Commands {
		count := 0
		if p.CommandUsage != nil {
			count = p.CommandUsage[c.Alias]
		}
		out = append(out, rankedCommand{Command: c, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Command.Alias < out[j].Command.Alias
	})
	return out
}

func (p *Project) TopCommands(n int) []Command {
	ranked := p.ranked()
	if n > len(ranked) {
		n = len(ranked)
	}
	top := make([]Command, 0, n)
	for i := 0; i < n; i++ {
		top = append(top, ranked[i].Command)
	}
	return top
}

func (p *Project) RestCommands() []Command {
	ranked := p.ranked()
	if len(ranked) <= TopCommandSlots {
		return nil
	}
	rest := make([]Command, 0, len(ranked)-TopCommandSlots)
	for i := TopCommandSlots; i < len(ranked); i++ {
		rest = append(rest, ranked[i].Command)
	}
	return rest
}
