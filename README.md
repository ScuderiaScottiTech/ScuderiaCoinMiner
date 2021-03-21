# ScuderiaCoinMiner

![](https://img.shields.io/github/go-mod/go-version/ScuderiaScottiTech/ScuderiaCoinMiner)
![](https://img.shields.io/github/v/tag/ScuderiaScottiTech/ScuderiaCoinMiner)
![](https://img.shields.io/github/workflow/status/ScuderiaScottiTech/ScuderiaCoinMiner/goreleaser?label=Releaser)

## Usage

By running the following command you'll be able to get a brief explainaion of all the parameters you can set on your miner

```bash
./ScuderiaCoinMiner --help
```
```
Usage of ./ScuderiaCoinMiner:
  -api string
        API endpoint of the mining server
  -goroutines int
        Number of goroutines to mine onto (default 1)
  -id string
        YOUR telegram ID
  -ratecounter
        Rate counter enabled (may degrade performance)
```

*__NOTE:__ Goroutines ARE NOT threads*

### Example usage

```bash
./ScuderiaCoinMiner -api https://mineapi.scuderia.tech -ratecounter -id <tuoid> -goroutines 6
```

Example output:
```
Getting new challenge info: Random: -3327122211037515956 Difficulty: 6 Reward: 10
Spawning goroutine
Spawning goroutine
Spawning goroutine
Spawning goroutine
Spawning goroutine
Spawning goroutine
Hash rate:  2113495 / second
Correct hash rate:  0 / hour
...
...
```
