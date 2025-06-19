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
|  | `duration`, `duration` | `duration` |
|  | `duration`, `timestamp` | `timestamp` |
|  | `timestamp`, `duration` | `timestamp` |
| `_-_` | `int`, `int` | `int` |
|  | `uint`, `uint` | `uint` |
|  | `double`, `double` | `double` |
|  | `duration`, `duration` | `duration` |
|  | `timestamp`, `duration` | `timestamp` |
|  | `timestamp`, `timestamp` | `duration` |
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
|  | `duration`, `duration` | `bool` |
|  | `timestamp`, `timestamp` | `bool` |
| `_<_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
|  | `duration`, `duration` | `bool` |
|  | `timestamp`, `timestamp` | `bool` |
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
|  | `duration`, `duration` | `bool` |
|  | `timestamp`, `timestamp` | `bool` |
| `_>_` | `int`, `int` | `bool` |
|  | `int`, `double` | `bool` |
|  | `int`, `uint` | `bool` |
|  | `uint`, `uint` | `bool` |
|  | `uint`, `int` | `bool` |
|  | `uint`, `double` | `bool` |
|  | `double`, `double` | `bool` |
|  | `double`, `int` | `bool` |
|  | `double`, `uint` | `bool` |
|  | `duration`, `duration` | `bool` |
|  | `timestamp`, `timestamp` | `bool` |
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
| `duration` | `duration` | `duration` |
|  | `int` | `duration` |
|  | `string` | `duration` |
| `endsWith` | `string`, `string` | `bool` |
| `get` | `map<dyn_A, dyn_B>`, `dyn_A` | `dyn_B` |
|  | `map<dyn_A, dyn_B>`, `dyn_A`, `dyn_B` | `dyn_B` |
| `getDate` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getDayOfMonth` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getDayOfWeek` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getDayOfYear` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getFullYear` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getHours` | `duration` | `int` |
|  | `timestamp` | `int` |
|  | `timestamp`, `string` | `int` |
| `getMilliseconds` | `duration` | `int` |
|  | `timestamp` | `int` |
|  | `timestamp`, `string` | `int` |
| `getMinutes` | `duration` | `int` |
|  | `timestamp` | `int` |
|  | `timestamp`, `string` | `int` |
| `getMonth` | `timestamp`, `string` | `int` |
|  | `timestamp` | `int` |
| `getSeconds` | `duration` | `int` |
|  | `timestamp` | `int` |
|  | `timestamp`, `string` | `int` |
| `has` | `map<dyn_A, dyn_B>`, `dyn_A` | `bool` |
| `indexOf` | `string`, `string`, `int` | `int` |
|  | `string`, `string` | `int` |
| `int` | `double` | `int` |
|  | `uint` | `int` |
|  | `int` | `int` |
|  | `string` | `int` |
|  | `duration` | `int` |
|  | `timestamp` | `int` |
| `join` | `list<string>`, `string` | `string` |
|  | `list<string>` | `string` |
| `lastIndexOf` | `string`, `string`, `int` | `int` |
|  | `string`, `string` | `int` |
| `lowerAscii` | `string` | `string` |
| `matches` | `string`, `string` | `bool` |
| `now` | - | `timestamp` |
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
|  | `duration` | `string` |
|  | `timestamp` | `string` |
| `substring` | `string`, `int`, `int` | `string` |
|  | `string`, `int` | `string` |
| `timestamp` | `timestamp` | `timestamp` |
|  | `int` | `timestamp` |
|  | `string` | `timestamp` |
| `trim` | `string` | `string` |
| `type` | `dyn_A` | `type` |
| `uint` | `double` | `uint` |
|  | `uint` | `uint` |
|  | `int` | `uint` |
|  | `string` | `uint` |
| `upperAscii` | `string` | `string` |
