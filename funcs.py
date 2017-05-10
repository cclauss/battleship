#!/usr/bin/env python3

from collections import defaultdict
import sys

# func_dict = {func_name: (start_line, end_line, lines)}

# def get_lines(filename='battleship.go'):
#    with open(filename) as in_file:
#        return [line.strip() for line in in_file]


def funcs(filename='battleship.go'):
    # pass one
    with open(filename) as in_file:
        lines = []  # all lines in in_file
        func_dict = {}  # key = function name, value = [start_line, end_line]
        curr_func = ''  # will contain the current function with a trailing '('
        for i, line in enumerate(in_file):
            line = line.replace('(h helloHandler) ', '')  # special case!!!
            lines.append(line)
            if line.lstrip().startswith('func'):
                if curr_func:
                    func_dict[curr_func].append(i - 1)  # record end_line
                print(i, curr_func)
                curr_func = line.partition(' ')[-1].partition('(')[0] + '('
                func_dict[curr_func] = [i + 1]  # record start_line
    print(sorted(func_dict))
    # assert False, 'Dude.'
    assert func_dict, 'No functions were found on the first pass!'
    func_dict[curr_func].append(i)  # record end_line
    # print(func_dict, '\n')
    # pass two
    print(len(lines), len(func_dict))
    for curr_func, start_end_lines in func_dict.items():
        # print(curr_func, start_end_lines)
        func_body = '\n'.join(lines[start_end_lines[0]:start_end_lines[1]])
        func_dict[curr_func] = []  # now record any functions that are called
        for called_func in func_dict:
            if called_func in func_body:
                func_dict[curr_func].append(called_func + ')')
    return {func + ')': tuple(called) for func, called in func_dict.items()}


print('=' * 25)
filename = sys.argv[1] if sys.argv[1:] else 'battleship.go'
d = funcs(filename)
for func in sorted(d):
    called = d[func]
    print('{}: {}: {}'.format('NODE' if called else 'LEAF', func, called))

print('\nReversing...')
r = defaultdict(list)
for key, values in d.items():
    for value in values:
        r[value].append(key)
for func in sorted(r):
    print('{} is called by {}'.format(func, r[func]))
