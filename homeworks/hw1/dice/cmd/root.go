package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tcarzverey/course-go-python/homeworks/hw1/dice/roller"
)

var (
	diceFormula   string
	verboseModeOn bool
	shortModeOn   bool
	helpModeOn    bool
	iterations    int
	seed          int64
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dice <dice-formula>",
	Short: "This application emulates a dice roller",
	Long: `This application emulates a dice roller. It uses a classic format NdM, where:
- N - is number of dices, that would be rolled simultaneously 
- M - is number of edges on the each dice. So each dice would give a random number in [1, N]

For N, M there is a bounds that utility would take: 1 <= N, M <= 100. Any other numbers would return an error code 1`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		diceFormula = args[0]

		re, _ := regexp.Compile(`(?m)^(\d{1,2}|100)d(\d{1,2}|100)$`)
		match := re.FindStringSubmatch(diceFormula)

		if len(match) == 0 {
			fmt.Println(fmt.Sprintf(`The current input {%s} doesn't match required constrains of the 
				dice formula. Use help to learn about format`, diceFormula))
			os.Exit(1)
		}

		dices, _ := strconv.Atoi(match[1])
		faces, _ := strconv.Atoi(match[2])

		req := roller.Request{
			NumberOfDices:      dices,
			NumberOfFaces:      faces,
			NumberOfIterations: iterations,
			Seed:               seed,
		}

		if err := req.Validate(); err != nil {
			fmt.Println(fmt.Sprintf("Current parameters of dices are not valid: %s", err))
			os.Exit(1)
		}

		if verboseModeOn && shortModeOn {
			fmt.Println("You can't activate both short print mode and verbose print mode")
			os.Exit(1)
		}

		res := roller.RollResult{
			Mode: roller.PrintModeDefault,
		}

		if verboseModeOn {
			res.Mode = roller.PrintModeVerbose
		}

		if shortModeOn {
			res.Mode = roller.PrintModeOnlySum
		}

		_ = res.Roll(&req)
		fmt.Print(res.ToString())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dice.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&helpModeOn, "help", "h", false, "Help message")

	rootCmd.Flags().BoolVarP(&verboseModeOn, "verbose", "v", false, `Using this command would made utility 
print value of each dice in each roll and also prints mean value of rolled dices in each iteration. This command 
can't be used with only sum mode(-s / --sum flag)'`)

	rootCmd.Flags().BoolVarP(&shortModeOn, "sum", "s", false, `Using this command would made
utility print only sum of the results of each iteration`)

	rootCmd.Flags().IntVarP(&iterations, "count", "n", 1, "Number of iterations to roll")

	rootCmd.Flags().Int64Var(&seed, "seed", time.Now().UnixNano(), "Seed for random number generator")
}
