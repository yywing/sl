# Functions

| Name | Input | Output |
|------|-------|--------|
| `!_` | `bool` | `bool` |
| `-_` | `int` | `int` |
|  | `double` | `double` |
| `_!=_` | `dyn_A`, `dyn_A` | `bool` |
| `_%_` | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
| `_&&_` | `bool`, `bool` | `bool` |
| `_*_` | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
|  | `double`, `double` | `double` |
| `_+_` | `bytes`, `bytes` | `bytes` |
|  | `double`, `double` | `double` |
|  | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
|  | `string`, `string` | `string` |
|  | `list<dyn_A>`, `list<dyn_A>` | `list<dyn_A>` |
| `_-_` | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
|  | `double`, `double` | `double` |
| `_/_` | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
|  | `double`, `double` | `double` |
| `_<=_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
| `_<_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
| `_==_` | `dyn_A`, `dyn_A` | `bool` |
| `_>=_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
| `_>_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
| `_in_` | `dyn_A`, `list<dyn_A>` | `bool` |
|  | `dyn_A`, `map<dyn_A, dyn_B>` | `bool` |
| `_\|\|_` | `bool`, `bool` | `bool` |
| `bool` | `bool` | `bool` |
|  | `string` | `bool` |
| `bytes` | `bytes` | `bytes` |
|  | `string` | `bytes` |
| `charAt` | `string`, `int` | `string` |
| `contains` | `string`, `string` | `bool` |
| `double` | `int` | `double` |
|  | `uint` | `double` |
|  | `double` | `double` |
|  | `string` | `double` |
| `endsWith` | `string`, `string` | `bool` |
| `get` | `map<dyn_A, dyn_B>`, `dyn_A` | `dyn_B` |
|  | `map<dyn_A, dyn_B>`, `dyn_A`, `dyn_B` | `dyn_B` |
| `has` | `map<dyn_A, dyn_B>`, `dyn_A` | `bool` |
| `indexOf` | `string`, `string`, `int` | `int` |
|  | `string`, `string` | `int` |
| `int` | `double` | `int` |
|  | `uint` | `int` |
|  | `int` | `int` |
|  | `string` | `int` |
| `join` | `list<string>`, `string` | `string` |
|  | `list<string>` | `string` |
| `lastIndexOf` | `string`, `string`, `int` | `int` |
|  | `string`, `string` | `int` |
| `lowerAscii` | `string` | `string` |
| `matches` | `string`, `string` | `bool` |
| `quote` | `string` | `string` |
| `replace` | `string`, `string`, `string`, `int` | `string` |
|  | `string`, `string`, `string` | `string` |
| `reverse` | `string` | `string` |
| `size` | `bytes` | `int` |
|  | `string` | `int` |
|  | `list<dyn_A>` | `int` |
|  | `map<dyn_A, dyn_B>` | `int` |
| `split` | `string`, `string`, `int` | `list<string>` |
|  | `string`, `string` | `list<string>` |
| `startsWith` | `string`, `string` | `bool` |
| `string` | `string` | `string` |
|  | `bytes` | `string` |
|  | `bool` | `string` |
|  | `double` | `string` |
|  | `int` | `string` |
|  | `uint` | `string` |
| `substring` | `string`, `int`, `int` | `string` |
|  | `string`, `int` | `string` |
| `trim` | `string` | `string` |
| `type` | `dyn_A` | `type` |
| `uint` | `double` | `uint` |
|  | `int` | `uint` |
|  | `uint` | `uint` |
|  | `string` | `uint` |
| `upperAscii` | `string` | `string` |
