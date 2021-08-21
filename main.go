/*
 * Script to generate NLU data to increase Classification and Extraction from Chatbot.

 * This script reads the file search/translations/french.csv and randomly choose some english phrases
 * contained in this file. After this, we format the phrase accondingly with templates defined in template variable.
 * After all, the result is saved in a new .txt file.
 */

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const maxLines = 50

var templates = []string{
	"can you tell me how i would normally say [{phrase}](phrase) as a [{lang}](lang) person",
	"can you tell me how i would say '[{phrase}](phrase)' in [{lang}](lang)",
	"can you tell me how to say '[{phrase}](phrase)' in [{lang}](lang)",
	"can you translate [{phrase}](phrase) into [{lang}](lang) for me",
	"could you translate [{phrase}](phrase) into [{lang}](lang)",
	"could you translate [{phrase}](phrase) into [{lang}](lang) for me",
	"do you know how to say [{phrase}](phrase) in [{lang}](lang)",
	"english to [{lang}](lang) for [{phrase}](phrase)",
	"how can i say [{phrase}](phrase) in [{lang}](lang)",
	"how can i [{phrase}](phrase] in [{lang}](lang)",
	"how could i say [{phrase}](phrase) in [{lang}](lang)",
	"how does one say [{phrase}](phrase) in [{lang}](lang)",
	"how do i say '[{phrase}](phrase)' in [{lang}](lang)",
	"how do i say [{phrase}](phrase) in [{lang}](lang)",
	"how do i say [{phrase}](phrase] in [{lang}](lang)",
	"how do i [{phrase}](phrase) in [{lang}](lang)",
	"how do [{lang}](lang) people say [{phrase}](phrase)",
	"how do [{lang}](lang) say [{phrase}](phrase)",
	"how do they say [{phrase}](phrase) in brazil](lang)",
	"how do they say [{phrase}](phrase) in [{lang}](lang)",
	"how do they say [{phrase}](phrase) in [{lang}](lang)",
	"how do you say [{phrase}](phrase) in [{lang}](lang)",
	"how do you say [{phrase}](phrase) in ]spanish](lang)",
	"how is [{phrase}](phrase) said in [{lang}](lang)",
	"how might i say [{phrase}](phrase) if i were [{lang}](lang)",
	"how should i say [{phrase}](phrase) in [{lang}](lang)",
	"how would i say if i were [{lang}](lang) [{phrase}](phrase)",
	"how would i say [{phrase}](phrase) if i were [{lang}](lang)",
	"how would i say '[{phrase}](phrase)' in [{lang}](lang)",
	"how would i say [{phrase}](phrase) in [{lang}](lang)",
	"how would one say [{phrase}](phrase) in [{lang}](lang)",
	"how would they say say [{phrase}](phrase) in [{lang}](lang)",
	"how would they say [{phrase}](phrase) in [{lang}](lang)",
	"how would you say [{phrase}](phrase) in [{lang}](lang)",
	"if i were [{lang}](lang) how would i say [{phrase}](phrase)",
	"if i were [{lang}](lang) how would i say that [{phrase}](phrase)",
	"i must know how to say [{phrase}](phrase) in [{lang}](lang)",
	"i need to know how to say [{phrase}](phrase) in [{lang}](lang)",
	"i need you to translate the sentence '[{phrase}](phrase)' into [{lang}](lang)",
	"in [{lang}](lang) how can I say [{phrase}](phrase)?",
	"in [{lang}](lang) how do i say [{phrase}](phrase)",
	"in [{lang}](lang) how do they say [{phrase}](phrase)",
	"in [{lang}](lang) [{phrase}](phrase) is said how",
	"i want to know how to say [{phrase}](phrase) in [{lang}](lang)",
	"i would i say [{phrase}](phrase) if i were [{lang}](lang)",
	"i would like to know how to say [{phrase}](phrase) in [{lang}](lang)",
	"i would like to know the proper way to [{phrase}](phrase) in [{lang}](lang)",
	"please tell me how to [{phrase}](phrase) in [{lang}](lang)",
	"please translate [{phrase}](phrase) into [{lang}](lang) for me",
	"[{phrase}](phrase) in [{lang}](lang)",
	"tell me how the [{lang}](lang) say [{phrase}](phrase)",
	"tell me how to say '[{phrase}](phrase)' in [{lang}](lang)",
	"tell me how to say [{phrase}](phrase) in [{lang}](lang)",
	"translate english to [{lang}](lang) [{phrase}](phrase)",
	"translate for me [{phrase}](phrase) into [{lang}](lang)",
	"translate [{phrase}](phrase) in [{lang}](lang)",
	"translate [{phrase}](phrase) into [{lang}](lang) for me",
	"translate [{phrase}](phrase) to [{lang}](lang)",
	"translate [{phrase}](phrase] to [{lang}](lang)",
	"what do i say for [{phrase}](phrase) in [{lang}](lang)",
	"what do [{lang}](lang) people say for the word [{phrase}](phrase)",
	"what do you [{phrase}](phrase) if you were [{lang}](lang)",
	"what expression would i use to say [{phrase}](phrase) if i were an [{lang}](lang)",
	"what is [{lang}](lang) for [{phrase}](phrase)",
	"what is [{phrase}](phrase) in [{lang}](lang)",
	"what is the correct way to say '[{phrase}](phrase)' in [{lang}](lang)",
	"what is the equivalent of '[{phrase}](phrase)' in [{lang}](lang)",
	"what is the right way to say [{phrase}](phrase) in [{lang}](lang)",
	"what is the way to say [{phrase}](phrase) in [{lang}](lang)",
	"what is the word for [{phrase}](phrase) [{lang}](lang)",
	"what phrase means [{phrase}](phrase) in [{lang}](lang)",
	"what [{lang}](lang) word means [{phrase}](phrase)",
	"what's local slang for [{phrase}](phrase) in [{lang}](lang)",
	"what's the [{lang}](lang) word for [{phrase}](phrase)",
	"what's the [{lang}](lang) word you use for [{phrase}](phrase)",
	"what's the word for [{phrase}](phrase) in [{lang}](lang)",
	"what words would i use to [{phrase}](phrase) if i were [{lang}](lang)",
	"what would the word for [{phrase}](phrase) be in [{lang}](lang)",
	"would you tell me how to say [{phrase}](phrase) in [{lang}](lang)",
}
var languages = []string{"portuguese", "french", "italian", "spanish", "chinese", "hindi", "arabic", "bengali", "russian", "japanese"}

func formatString(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func normalizeString(s string) string {
	s = strings.ToLower(s)
	re, _ := regexp.Compile(`[^\w]`)
	s = re.ReplaceAllString(s, "")

	return s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	csvFile, err := os.Open("./search/translations/french.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println("error reading csv", err)
		return
	}

	count := 0
	generatedPhrases := []string{}
	for _, line := range csvLines {
		if count >= maxLines {
			break
		}

		if rand.Intn(2) == 1 {
			count++
			englishSentence := normalizeString(line[0])
			for _, template := range templates {
				language := languages[rand.Intn(10)]

				formattedString := formatString("- "+template, "{phrase}", englishSentence, "{lang}", language)
				formattedString = strings.ToLower(formattedString)

				generatedPhrases = append(generatedPhrases, formattedString)
			}
		}
	}

	saveResultIntoFile(generatedPhrases)
}

func saveResultIntoFile(result []string) {
	file, err := os.Create("data.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)

	for _, data := range result {
		datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}
