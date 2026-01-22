#!/usr/bin/env python
# coding: utf-8

# ## Graded Lab: Reflection in a Research Agent
# 
# In this graded lab, you’ll implement a simple **agentic workflow** designed to simulate reflective thinking in a writing task. This is one building block of a more complex research agent that will be constructed throughout the course.
# 
# ### Objective
# 
# Build a three-step workflow where an LLM writes an essay draft, critiques it, and rewrites it. 
# 
# * **Step 1 – Drafting:** Call the LLM to generate an initial draft of an essay based on a simple prompt.
# * **Step 2 – Reflection:** Reflect on the draft using a reasoning step. (Optionally, this can be done with a different model.)
# * **Step 3 – Revision:** Apply the feedback from the reflection to generate a revised version of the essay.
# 

# ---
# <a name='submission'></a>
# 
# <h4 style="color:green; font-weight:bold;">TIPS FOR SUCCESSFUL GRADING OF YOUR ASSIGNMENT:</h4>
# 
# * All cells are frozen except for the ones where you need to write your solution code or when explicitly mentioned you can interact with it.
# 
# * In each exercise cell, look for comments `### START CODE HERE ###` and `### END CODE HERE ###`. These show you where to write the solution code. **Do not add or change any code that is outside these comments**.
# 
# * You can add new cells to experiment but these will be omitted by the grader, so don't rely on newly created cells to host your solution code, use the provided places for this.
# 
# * Avoid using global variables unless you absolutely have to. The grader tests your code in an isolated environment without running all cells from the top. As a result, global variables may be unavailable when scoring your submission. Global variables that are meant to be used will be defined in UPPERCASE.
# 
# * To submit your notebook for grading, first save it by clicking the 💾 icon on the top left of the page and then click on the <span style="background-color: red; color: white; padding: 3px 5px; font-size: 16px; border-radius: 5px;">Submit assignment</span> button on the top right of the page.
# ---

# Before interacting with the language models, we initialize the `aisuite` client. This setup loads environment variables (e.g., API keys) from a `.env` file to securely authenticate with backend services. The `ai.Client()` instance will be used to make all model calls throughout this workflow.

# In[1]:


from dotenv import load_dotenv

load_dotenv()

import aisuite as ai

# Define the client. You can use this variable inside your graded functions!
CLIENT = ai.Client()


# In[2]:


import unittests


# ## Exercise 1: `generate_draft` Function
# 
# **Objective**:
# Write a function called `generate_draft` that takes in a string topic and uses a language model to generate a complete draft essay.
# 
# **Inputs**:
# 
# * `topic` (str): The essay topic.
# * `model` (str, optional): The model identifier to use. Defaults to `"openai:gpt-4o"`.
# 
# **Output**:
# 
# * A string representing the full draft of the essay.
# 
# The setup for calling the LLM using the aisuite library is already provided. Focus on crafting the prompt content. You can reference this setup in later exercises to understand how to interact with the library effectively.
# 

# In[9]:


# GRADED FUNCTION: generate_draft

def generate_draft(topic: str, model: str = "openai:gpt-4o") -> str: 
    
    ### START CODE HERE ###

    # Define your prompt here. A multi-line f-string is typically used for this.
    prompt = f"""
    You are an essay writer given a topic you will write an essay for that topic
    
    Topic:
    {topic}
    
    Respond with only the essay without any other words before or after the essay
    """

    ### END CODE HERE ###
    
    # Get a response from the LLM by creating a chat with the client.
    response = CLIENT.chat.completions.create(
        model=model,
        messages=[{"role": "user", "content": prompt}],
        temperature=1.0,
    )

    return response.choices[0].message.content


# Run the following cell to check your code is working correctly:

# In[11]:


# Test your code!
unittests.test_generate_draft(generate_draft)


# ## Exercise 2: `reflect_on_draft` Function
# 
# **Objective**:
# Write a function called `reflect_on_draft` that takes a previously generated essay draft and uses a language model to provide constructive feedback.
# 
# **Inputs**:
# 
# * `draft` (str): The essay text to reflect on.
# * `model` (str, optional): The model identifier to use. Defaults to `"openai:o4-mini"`.
# 
# **Output**:
# 
# * A string with feedback in paragraph form.
# 
# **Requirements**:
# 
# * The feedback should be critical but constructive.
# * It should address issues such as structure, clarity, strength of argument, and writing style.
# * The function should send the draft to the model and return its response.
# 
# You do **not** need to rewrite the essay at this step—just analyze and reflect on it.
# 

# In[14]:


# GRADED FUNCTION: reflect_on_draft

def reflect_on_draft(draft: str, model: str = "openai:o4-mini") -> str:

    ### START CODE HERE ###

    # Define your prompt here. A multi-line f-string is typically used for this.
    prompt = f"""
    You are a professional editor. Given an essay write a your feedback about this essay.
    The feedback should be critical but constructive.
    It should address issues such as structure, clarity, strength of argument, and writing style.
    
    Essay:
    {draft}
    
    Respond with only the feedback without any words before or after the feedback
    
    """

    ### END CODE HERE ###

    # Get a response from the LLM by creating a chat with the client.
    response = CLIENT.chat.completions.create(
        model=model,
        messages=[{"role": "user", "content": prompt}],
        temperature=1.0,
    )

    return response.choices[0].message.content


# In[15]:


# Test your code!
unittests.test_reflect_on_draft(reflect_on_draft)


# ## Exercise 3: `revise_draft` Function
# 
# **Objective**:
# Implement a function called `revise_draft` that improves a given essay draft based on feedback from a reflection step.
# 
# **Inputs**:
# 
# * `original_draft` (str): The initial version of the essay.
# * `reflection` (str): Constructive feedback or critique on the draft.
# * `model` (str, optional): The model identifier to use. Defaults to `"openai:gpt-4o"`.
# 
# **Output**:
# 
# * A string containing the revised and improved essay.
# 
# **Requirements**:
# 
# * The revised draft should address the issues mentioned in the feedback.
# * It should improve clarity, coherence, argument strength, and overall flow.
# * The function should use the feedback to guide the revision, and return only the final revised essay.
# 
# In this final exercise, you'll also need to manage the call to the LLM using the CLIENT, as you've practiced in previous exercises.

# In[24]:


# GRADED FUNCTION: revise_draft

def revise_draft(original_draft: str, reflection: str, model: str = "openai:gpt-4o") -> str:

    ### START CODE HERE ###

    # Define your prompt here. A multi-line f-string is typically used for this.
    prompt = f"""
    You are a professional writer. You are given an initial draft on an essay and a feedback on that essay.
    You should rewrite the essay.
    The revised draft should address the issues mentioned in the feedback.
    It should improve clarity, coherence, argument strength, and overall flow.
    
    Essay:
    {original_draft}
    
    Feedback:
    {reflection}
    
    Respond with only the revised essay without any words before or after it
    
    """

    # Get a response from the LLM by creating a chat with the client.
    response = CLIENT.chat.completions.create(
        model = model,
        messages = [{"role":"user","content":prompt}],
        temperature=1.0,
    )

    ### END CODE HERE ###

    return response.choices[0].message.content


# In[25]:


# Test your code!
unittests.test_revise_draft(revise_draft)


# ### 🧪 Test the Reflective Writing Workflow
# 
# Use the functions you implemented to simulate the complete writing workflow:
# 
# 1. **Generate a draft** in response to the essay prompt.
# 2. **Reflect** on the draft to identify improvements.
# 3. **Revise** the draft using the feedback.
# 
# Observe the outputs of each step. You do **not** need to modify the outputs — just verify that the workflow runs as expected and each component returns a valid string.
# 

# In[ ]:


essay_prompt = "Should social media platforms be regulated by the government?"

# Agent 1 – Draft
draft = generate_draft(essay_prompt)
print("📝 Draft:\n")
print(draft)

# Agent 2 – Reflection
feedback = reflect_on_draft(draft)
print("\n🧠 Feedback:\n")
print(feedback)

# Agent 3 – Revision
revised = revise_draft(draft, feedback)
print("\n✍️ Revised:\n")
print(revised)


# To better visualize the output of each step in the reflective writing workflow, we use a utility function called `show_output`. This function displays the results of each stage (drafting, reflection, and revision) in styled boxes with custom background and text colors, making it easier to compare and understand the progression of the essay.
# 

# In[ ]:


from utils import show_output

essay_prompt = "Should social media platforms be regulated by the government?"

show_output("Step 1 – Draft", draft, background="#fff8dc", text_color="#333333")
show_output("Step 2 – Reflection", feedback, background="#e0f7fa", text_color="#222222")
show_output("Step 3 – Revision", revised, background="#f3e5f5", text_color="#222222")


# ## Check grading feedback
# 
# If you have collapsed the right panel to have more screen space for your code, as shown below:
# 
# <img src="./images/collapsed.png" alt="Collapsed Image" width="800" height="400"/>
# 
# You can click on the left-facing arrow button (highlighted in red) to view feedback for your submission after submitting it for grading. Once expanded, it should display like this:
# 
# <img src="./images/expanded.png" alt="Expanded Image" width="800" height="400"/>
