# Function Deployment
Functions can be effortlessly deployed through the Nuclio graphical user interface (GUI). To ensure the seamless operation of the workflow and to avoid a detour through the nuclio controller, these functions have to publish their output to the respective queue.

## Example Workflow
In this directory are three go functions that form a workflow, which first converts an input string to uppercase, then reverts it and finally stretches it. (Example: nuclio => O I L C U N)
Each function publishes its output to the corresponding queue of the next function in the worklow, which triggers the corresponding grainworker to take the message and invoke the next function.