#!/home/nox/micromamba/envs/shell-helper/bin/python

import sys
import re
import openai
import subprocess
import os

API_KEYFILE = os.path.join(os.path.expanduser("~"), ".config", "shell-helper", "key")

# Set up your OpenAI API key
with open(API_KEYFILE, 'r') as file:
    openai.api_key = file.read().strip()

def get_chatgpt_response(prompt):
    response = openai.ChatCompletion.create(
        model="gpt-3.5-turbo",
        messages=[
            {"role": "system", "content": "You write modern, excellent concise bash oneliners. All codeblocks must start with the appropriate language identifiers."},
            {"role": "user", "content": prompt}
        ]
    )
    return response.choices[0].message['content'].strip()

def extract_code_block(text):
    pattern = r"```bash\n(.*?)\n```"
    match = re.search(pattern, text, re.DOTALL)
    if match:
        return match.group(1).strip()
    return None

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: ? <prompt>")
        sys.exit(1)

    user_input = " ".join(sys.argv[1:])
    prompt = f"I have the linux terminal open, I want to {user_input}. Output the correct bash one-liner to do this in a single codeblock. No japping!"
    
    response = get_chatgpt_response(prompt)
    print(response)

    shell_cmd = extract_code_block(response)
    if shell_cmd:
        user_input = input("Shell command detected. Execute?")
        if user_input.lower() == '':
            subprocess.run(shell_cmd, shell=True)