from typing import Any, Text, Dict, List

from rasa_sdk import Action, Tracker
from rasa_sdk.executor import CollectingDispatcher

from actions.utils import isEmptyTextOrNone
from actions.elasticsearch import save_translation, findTranslationsInElasticsearch
from actions.google import translate_text, load_languages_in_memory

load_languages_in_memory()

class ActionTranslate(Action):

    def name(self) -> Text:
        return "action_translate"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        slots = tracker.current_slot_values()
        
        sentence = slots["phrase"]
        language = slots["lang"]

        if isEmptyTextOrNone(sentence) or isEmptyTextOrNone(language):
            dispatcher.utter_message(text="Sorry I can't understand what you to translate")
            return []

        # Search matches in ES
        response = findTranslationsInElasticsearch(sentence, language)

        # Display result to the user
        if len(response.hits) > 0: 
            dispatcher.utter_message(text=u"The translation is: {}".format(response.hits[0]["sentence"]))
        else:
            dispatcher.utter_message(text=u"I didn't find anything in my database. So I'll search on Google")
            # If don't find in ES, translate in Google and save it to ES
            response = ""
            try:
                response = translate_text(sentence, language)
            except KeyError:
                dispatcher.utter_message(text="Sorry I can't translate to this language.")
                return []

            # Display result to the user
            dispatcher.utter_message(text=u"The translation is: {}".format(response))

            # Save result in ES so the search can have more accuracy
            save_translation(sentence, response, language)

        return []

class ActionTranslateGoogle(Action):

    def name(self) -> Text:
        return "action_translate_google"

    def run(self, dispatcher: CollectingDispatcher,
            tracker: Tracker,
            domain: Dict[Text, Any]) -> List[Dict[Text, Any]]:

        slots = tracker.current_slot_values()
        sentence = slots["phrase"]
        language = slots["lang"]

        if isEmptyTextOrNone(sentence) or isEmptyTextOrNone(language):
            dispatcher.utter_message(text="Sorry I can't understand what you to translate")
            return []

        # Translate text using Google
        response = ""
        try:
            response = translate_text(sentence, language)
        except KeyError:
            dispatcher.utter_message(text="Sorry I can't translate to this language.")
            return []

        # Display result to the user
        dispatcher.utter_message(text=u"The translation is: {}".format(response))

        # Save result in ES so the search can have more accuracy
        save_translation(sentence, response, language)

        return []