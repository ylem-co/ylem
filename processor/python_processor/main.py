from typing import Union
from fastapi import FastAPI, Response, status
import json
import re
from pydantic import BaseModel
import logging
import time
import datetime
from RestrictedPython import compile_restricted, safe_globals
from RestrictedPython import safe_builtins
from RestrictedPython import limited_builtins
from RestrictedPython import utility_builtins
from RestrictedPython.Eval import default_guarded_getitem, default_guarded_getiter
from AccessControl.ZopeGuards import protected_inplacevar
import pandas


class RequestBody(BaseModel):
    code: str
    input: str

safe_builtins.update(utility_builtins)
safe_builtins.update(limited_builtins)
safe_builtins['pandas'] = pandas
safe_builtins['json'] = json
safe_builtins['re'] = re
safe_builtins['time'] = time
safe_builtins['datetime'] = datetime

app = FastAPI()

@app.post("/eval")
def eval(body: RequestBody, response: Response):
    _getitem_ = default_guarded_getitem
    _getiter_ = default_guarded_getiter
    _inplacevar_ = protected_inplacevar
    
    code = body.code
    input = json.loads(body.input)

    code = re.sub("""(?<!")import +\n*socket *(?!")""", "", code)
    code = re.sub("""(?<!")import +\n*requests *(?!")""", "", code)

    try:
        ldict = locals()
        byte_code = compile_restricted(code, filename='<inline code>', mode='exec')
        exec(byte_code, {'__builtins__': safe_builtins}, ldict)
        
        return ldict["input"]
    except SyntaxError as err:
        response.status_code = 400
        return json.dumps(err.msg)
    except NameError as err:
        response.status_code = 500
        return json.dumps(err.args[0])
    except Exception as err:
        response.status_code = 500
        return json.dumps(err.__str__())
