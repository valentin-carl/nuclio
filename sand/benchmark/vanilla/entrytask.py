from nuclio_sdk import Event

def handler(context, event):

    target, actual = int(event.body), 0
    while actual < target:
        res = context.platform.call_function('additiontask-vanilla', Event(body=f"{actual}", method="GET"))
        actual = int(res.body)
    
    return actual
