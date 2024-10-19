# YLEM PYTHON CODE PROCESSOR

<a href="https://github.com/ylem-co/ylem?tab=Apache-2.0-1-ov-file">![Static Badge](https://img.shields.io/badge/license-Apache%202.0-black)</a>
<a href="https://ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/website-ylem.co-black)</a>
<a href="https://docs.ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/documentation-docs.ylem.co-black)</a>
<a href="https://join.slack.com/t/ylem-co/shared_invite/zt-2nawzl6h0-qqJ0j7Vx_AEHfnB45xJg2Q" target="_blank">![Static Badge](https://img.shields.io/badge/community-join%20Slack-black)</a>

Python code processor is an API for evaluating Python expressions from the "Code" pipeline task. It executes arbitrary Python code and returns the results.

It is available inside the Ylem network on http://ylem_python_processor:7338 or from the host machine on http://127.0.0.1:7338.

# Endpoints

## POST /eval

### Request body:

```js
{
    "code": "input['value'] = 2", // the code to execute
    "input": "{\"value\": 1}" // input value, available in the code as "input" variable
}
```

### Response body:
```js
{
    "statusCode": 200,
    "body": "{\"value\": 2}" // execution result
}
```
