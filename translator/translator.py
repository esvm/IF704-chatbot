from google.cloud import translate_v2 as translate
import six
import os
os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = os.path.abspath(os.getcwd() + "\\..\\credenciais.json")
def list_languages():
    translate_client = translate.Client()

    results = translate_client.get_languages()

    for language in results:
        print(u"{name} ({language})".format(**language))

linguas = input("Linguas disponiveis? [Y/n]: ")
while (linguas != "Y") & (linguas != "n"):
    linguas = input("Linguas disponiveis? [Y/n]: ")
if linguas == "Y":
    list_languages()

def translate_text(target, text):
    translate_client = translate.Client()

    if isinstance(text, six.binary_type):
        text = text.decode("utf-8")
    result = translate_client.translate(text, target_language=target)


    print(u"Tradução: {}".format(result["translatedText"]))

target = input("Digite a lingua alvo: ")
text = input("Digite a frase a ser traduzida: ")
translate_text(target, text)