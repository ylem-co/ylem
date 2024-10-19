import axios from "axios";

const STATS_API_URL = "/stats-api/";

class StatService {
    getPipelineStats(uuid, dateFrom, dateTo) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/stats/' + dateFrom + '/' + dateTo,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {
            data: {
                "num_of_successes": 3567,
                "num_of_failures": 98,
                "average_duration": 5678,
                "is_last_run_successful": false,
                "last_run_duration": 12345,
                "last_run_executed_at": "2021-01-05 05:27:00"
            }
        };*/
    }

    getLastRunsLog(uuid, dateFrom, dateTo) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/last-runs-log/' + dateFrom + '/' + dateTo,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
        /*return {"data": [
    {
        "pipeline_run_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 13500,
        "executed_at": "2024-01-13T18:36:58Z"
    },
    {
        "pipeline_run_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 639,
        "output": "",
        "metric_value": 0,
        "executed_at": "2024-01-13T18:36:57Z"
    },
    {
        "pipeline_run_uuid": "7ddae87a-4ff4-4fa0-8bd9-a88ec21f48a3",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 583,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:32:57Z"
    },
    {
        "pipeline_run_uuid": "7ddae87a-4ff4-4fa0-8bd9-a88ec21f48a3",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 13500,
        "executed_at": "2023-07-29T12:32:57Z"
    },
    {
        "pipeline_run_uuid": "9aa7e69f-4b94-4c07-a230-1f56a5dc4b42",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 8500,
        "executed_at": "2023-07-29T12:32:21Z"
    },
    {
        "pipeline_run_uuid": "9aa7e69f-4b94-4c07-a230-1f56a5dc4b42",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 582,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:32:20Z"
    },
    {
        "pipeline_run_uuid": "18fb9f26-9603-4586-80d8-55c021f04cbc",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 567,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:31:37Z"
    },
    {
        "pipeline_run_uuid": "18fb9f26-9603-4586-80d8-55c021f04cbc",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 9700,
        "executed_at": "2023-07-29T12:31:37Z"
    },
    {
        "pipeline_run_uuid": "81e44070-04cb-4824-842f-345820fb3f03",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 10500,
        "executed_at": "2023-07-29T12:31:15Z"
    },
    {
        "pipeline_run_uuid": "81e44070-04cb-4824-842f-345820fb3f03",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 555,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:31:14Z"
    },
    {
        "pipeline_run_uuid": "6203cb69-3cf6-432e-8f00-da1df1b45f29",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 560,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:30:53Z"
    },
    {
        "pipeline_run_uuid": "6203cb69-3cf6-432e-8f00-da1df1b45f29",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 10000,
        "executed_at": "2023-07-29T12:30:53Z"
    },
    {
        "pipeline_run_uuid": "47268725-3bb3-456a-b662-b7fcd80af4b4",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 1,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T12:15:03Z"
    },
    {
        "pipeline_run_uuid": "47268725-3bb3-456a-b662-b7fcd80af4b4",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 567,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T12:15:02Z"
    },
    {
        "pipeline_run_uuid": "c6dd909f-7e21-4fb6-b404-ec4e2b3564da",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T11:15:02Z"
    },
    {
        "pipeline_run_uuid": "c6dd909f-7e21-4fb6-b404-ec4e2b3564da",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 581,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T11:15:01Z"
    },
    {
        "pipeline_run_uuid": "d6341e65-7b7e-407e-86ca-d55c3b408daa",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T10:15:02Z"
    },
    {
        "pipeline_run_uuid": "d6341e65-7b7e-407e-86ca-d55c3b408daa",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 611,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T10:15:01Z"
    },
    {
        "pipeline_run_uuid": "10b53679-5acf-43f0-a0b1-73dd4623da3c",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T09:15:02Z"
    },
    {
        "pipeline_run_uuid": "10b53679-5acf-43f0-a0b1-73dd4623da3c",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 585,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T09:15:01Z"
    },
    {
        "pipeline_run_uuid": "220e5289-2450-42bf-819a-d0f08da15467",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 1,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T08:15:02Z"
    },
    {
        "pipeline_run_uuid": "220e5289-2450-42bf-819a-d0f08da15467",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 570,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T08:15:01Z"
    },
    {
        "pipeline_run_uuid": "d03709dc-124b-4e57-a3f3-ff6e60c25bb7",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T07:15:02Z"
    },
    {
        "pipeline_run_uuid": "d03709dc-124b-4e57-a3f3-ff6e60c25bb7",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 644,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T07:15:01Z"
    },
    {
        "pipeline_run_uuid": "f5101287-7fde-4b6c-9b0b-f1a92efcf052",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T06:15:02Z"
    },
    {
        "pipeline_run_uuid": "f5101287-7fde-4b6c-9b0b-f1a92efcf052",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 587,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T06:15:01Z"
    },
    {
        "pipeline_run_uuid": "f3fdffdb-479c-4886-871d-123de42c58e8",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T05:15:02Z"
    },
    {
        "pipeline_run_uuid": "f3fdffdb-479c-4886-871d-123de42c58e8",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 569,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T05:15:01Z"
    },
    {
        "pipeline_run_uuid": "abf51d68-5027-4bce-8f40-099f98f48bec",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T04:15:02Z"
    },
    {
        "pipeline_run_uuid": "abf51d68-5027-4bce-8f40-099f98f48bec",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 585,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T04:15:01Z"
    },
    {
        "pipeline_run_uuid": "11346c60-ce23-49b7-9a79-e0e577a776c9",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T03:15:02Z"
    },
    {
        "pipeline_run_uuid": "11346c60-ce23-49b7-9a79-e0e577a776c9",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 914,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T03:15:01Z"
    },
    {
        "pipeline_run_uuid": "598a723c-865d-43b5-934a-ac7a81a2cea0",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T02:15:02Z"
    },
    {
        "pipeline_run_uuid": "598a723c-865d-43b5-934a-ac7a81a2cea0",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 572,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T02:15:01Z"
    },
    {
        "pipeline_run_uuid": "58f0666f-d4f6-418e-9ff7-d6da79a5f475",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T01:15:02Z"
    },
    {
        "pipeline_run_uuid": "58f0666f-d4f6-418e-9ff7-d6da79a5f475",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 597,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T01:15:01Z"
    },
    {
        "pipeline_run_uuid": "8c81ed22-dcc1-441e-ae10-4f2830ac1355",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-29T00:15:02Z"
    },
    {
        "pipeline_run_uuid": "8c81ed22-dcc1-441e-ae10-4f2830ac1355",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 583,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-29T00:15:01Z"
    },
    {
        "pipeline_run_uuid": "af88f0a6-264f-4928-acf2-87031125573f",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T23:15:02Z"
    },
    {
        "pipeline_run_uuid": "af88f0a6-264f-4928-acf2-87031125573f",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 574,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T23:15:01Z"
    },
    {
        "pipeline_run_uuid": "f39b92df-2438-4ece-b30b-54555d57677e",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T22:15:02Z"
    },
    {
        "pipeline_run_uuid": "f39b92df-2438-4ece-b30b-54555d57677e",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 607,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T22:15:01Z"
    },
    {
        "pipeline_run_uuid": "22a16a8c-bcf2-44ab-b9e5-33dfaa1b2bf3",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T21:15:02Z"
    },
    {
        "pipeline_run_uuid": "22a16a8c-bcf2-44ab-b9e5-33dfaa1b2bf3",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 574,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T21:15:01Z"
    },
    {
        "pipeline_run_uuid": "619a1e10-eea8-46a3-95a6-f6fcfe7e3675",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T20:15:03Z"
    },
    {
        "pipeline_run_uuid": "619a1e10-eea8-46a3-95a6-f6fcfe7e3675",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 678,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T20:15:02Z"
    },
    {
        "pipeline_run_uuid": "467efe74-5d52-4554-b899-597860bb2a79",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T19:15:02Z"
    },
    {
        "pipeline_run_uuid": "467efe74-5d52-4554-b899-597860bb2a79",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 673,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T19:15:02Z"
    },
    {
        "pipeline_run_uuid": "56a2a983-eca2-49f5-b03f-b1578c2ec7fd",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T18:15:02Z"
    },
    {
        "pipeline_run_uuid": "56a2a983-eca2-49f5-b03f-b1578c2ec7fd",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 678,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T18:15:01Z"
    },
    {
        "pipeline_run_uuid": "e76ddda0-2f70-4f83-9c4a-eafc947afb16",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 566,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T18:00:42Z"
    },
    {
        "pipeline_run_uuid": "e76ddda0-2f70-4f83-9c4a-eafc947afb16",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T18:00:42Z"
    },
    {
        "pipeline_run_uuid": "65df8e3c-e4e8-4cc9-9740-284b259d9e7a",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T18:00:35Z"
    },
    {
        "pipeline_run_uuid": "65df8e3c-e4e8-4cc9-9740-284b259d9e7a",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 654,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T18:00:34Z"
    },
    {
        "pipeline_run_uuid": "8c55cc0c-e0e8-4c54-b0d1-a1e574935d88",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1481.454,
        "executed_at": "2023-07-28T18:00:13Z"
    },
    {
        "pipeline_run_uuid": "8c55cc0c-e0e8-4c54-b0d1-a1e574935d88",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 658,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T18:00:13Z"
    },
    {
        "pipeline_run_uuid": "d18ecfef-aa05-4ca6-8eb1-f46c8d4ba216",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 722,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T17:59:00Z"
    },
    {
        "pipeline_run_uuid": "d18ecfef-aa05-4ca6-8eb1-f46c8d4ba216",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 961.454,
        "executed_at": "2023-07-28T17:59:00Z"
    },
    {
        "pipeline_run_uuid": "6250d2c3-e9b1-4388-b9fa-59393161d21b",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 961.454,
        "executed_at": "2023-07-28T17:58:44Z"
    },
    {
        "pipeline_run_uuid": "6250d2c3-e9b1-4388-b9fa-59393161d21b",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 638,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T17:58:43Z"
    },
    {
        "pipeline_run_uuid": "16ce28da-8c42-4cc8-b7f5-0da3164e6dc6",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 664,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T17:58:19Z"
    },
    {
        "pipeline_run_uuid": "16ce28da-8c42-4cc8-b7f5-0da3164e6dc6",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1848.454,
        "executed_at": "2023-07-28T17:58:19Z"
    },
    {
        "pipeline_run_uuid": "e0b90d63-ade0-4b1c-b7db-b9d210fa5fe0",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1817.454,
        "executed_at": "2023-07-28T17:57:42Z"
    },
    {
        "pipeline_run_uuid": "e0b90d63-ade0-4b1c-b7db-b9d210fa5fe0",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 675,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T17:57:41Z"
    },
    {
        "pipeline_run_uuid": "1551fd32-5ca2-4f4c-b3c0-bb46ff567952",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T17:15:02Z"
    },
    {
        "pipeline_run_uuid": "1551fd32-5ca2-4f4c-b3c0-bb46ff567952",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 669,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T17:15:01Z"
    },
    {
        "pipeline_run_uuid": "758a3f42-2ccf-4570-9b13-fe34b62cc4a3",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T16:15:02Z"
    },
    {
        "pipeline_run_uuid": "758a3f42-2ccf-4570-9b13-fe34b62cc4a3",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 676,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T16:15:01Z"
    },
    {
        "pipeline_run_uuid": "2dd7b00b-d745-4220-b13a-166bcf2240e5",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 669,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T15:15:01Z"
    },
    {
        "pipeline_run_uuid": "2dd7b00b-d745-4220-b13a-166bcf2240e5",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T15:15:01Z"
    },
    {
        "pipeline_run_uuid": "8b9eed3d-eee6-47ce-ac81-7e36787414e6",
        "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
        "task_type": "aggregator",
        "is_successful": true,
        "duration": 0,
        "output": "",
        "metric_value": 1461.454,
        "executed_at": "2023-07-28T14:45:54Z"
    },
    {
        "pipeline_run_uuid": "8b9eed3d-eee6-47ce-ac81-7e36787414e6",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": true,
        "duration": 945,
        "output": "",
        "metric_value": 0,
        "executed_at": "2023-07-28T14:45:53Z"
    },
    {
        "pipeline_run_uuid": "e2e0a319-c118-44c4-bcec-e3f63a0d5d6c",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": false,
        "duration": 3024,
        "output": "W3siYW1vdW50IjoyOTkuNzYsImNyZWF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjUzOjM1IiwiY3VycmVuY3kiOiJFVVIiLCJjdXN0b21lcl9mcm9tX2lkIjoxLCJjdXN0b21lcl90b19pZCI6MiwiaWQiOjEsInVwZGF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjUzOjM1In0seyJhbW91bnQiOjExLjAxLCJjcmVhdGVkX2F0IjoiMjAyMi0wOS0zMCAxODo1NzozNSIsImN1cnJlbmN5IjoiRVVSIiwiY3VzdG9tZXJfZnJvbV9pZCI6MSwiY3VzdG9tZXJfdG9faWQiOjIsImlkIjoyLCJ1cGRhdGVkX2F0IjoiMjAyMi0wOS0zMCAxODo1MzozNSJ9LHsiYW1vdW50IjoxMS4wMSwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTM6MzUiLCJjdXJyZW5jeSI6IkVVUiIsImN1c3RvbWVyX2Zyb21faWQiOjEsImN1c3RvbWVyX3RvX2lkIjoyLCJpZCI6MywidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTM6MzUifSx7ImFtb3VudCI6MTMwMCwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MjYiLCJjdXJyZW5jeSI6IlVTRCIsImN1c3RvbWVyX2Zyb21faWQiOjMsImN1c3RvbWVyX3RvX2lkIjo1LCJpZCI6NCwidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MjYifSx7ImFtb3VudCI6MTMwMCwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MjYiLCJjdXJyZW5jeSI6IlVTRCIsImN1c3RvbWVyX2Zyb21faWQiOjMsImN1c3RvbWVyX3RvX2lkIjo1LCJpZCI6NSwidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MjYifSx7ImFtb3VudCI6MTUwNC4zNiwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MzkiLCJjdXJyZW5jeSI6IlVTRCIsImN1c3RvbWVyX2Zyb21faWQiOjUsImN1c3RvbWVyX3RvX2lkIjozLCJpZCI6NiwidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTQ6MzkifSx7ImFtb3VudCI6MzQuODksImNyZWF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjU1OjEwIiwiY3VycmVuY3kiOiJFVVIiLCJjdXN0b21lcl9mcm9tX2lkIjo0LCJjdXN0b21lcl90b19pZCI6MSwiaWQiOjcsInVwZGF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjU1OjEwIn0seyJhbW91bnQiOjU2NC44OSwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTU6MjgiLCJjdXJyZW5jeSI6IkVVUiIsImN1c3RvbWVyX2Zyb21faWQiOjEsImN1c3RvbWVyX3RvX2lkIjo0LCJpZCI6OCwidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTU6MjgifSx7ImFtb3VudCI6MTU2Ny41MSwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTU6MjgiLCJjdXJyZW5jeSI6IkVVUiIsImN1c3RvbWVyX2Zyb21faWQiOjEsImN1c3RvbWVyX3RvX2lkIjo0LCJpZCI6OSwidXBkYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTU6MjgifSx7ImFtb3VudCI6MjEuMTEsImNyZWF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjU2OjA2IiwiY3VycmVuY3kiOiJFVVIiLCJjdXN0b21lcl9mcm9tX2lkIjoyLCJjdXN0b21lcl90b19pZCI6MSwiaWQiOjEwLCJ1cGRhdGVkX2F0IjoiMjAyMi0wOS0zMCAxODo1NjowNiJ9LHsiYW1vdW50IjozNC4zNSwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTY6MjYiLCJjdXJyZW5jeSI6IlVTRCIsImN1c3RvbWVyX2Zyb21faWQiOjMsImN1c3RvbWVyX3RvX2lkIjo1LCJpZCI6MTEsInVwZGF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjU2OjI2In0seyJhbW91bnQiOjExMi41NiwiY3JlYXRlZF9hdCI6IjIwMjItMDktMzAgMTg6NTY6NDQiLCJjdXJyZW5jeSI6IlVTRCIsImN1c3RvbWVyX2Zyb21faWQiOjUsImN1c3RvbWVyX3RvX2lkIjozLCJpZCI6MTIsInVwZGF0ZWRfYXQiOiIyMDIyLTA5LTMwIDE4OjU2OjQ0In1d",
        "metric_value": 0,
        "executed_at": "2023-07-28T14:45:40Z"
    },
    {
        "pipeline_run_uuid": "a1cbc9b9-fc34-4e56-91c4-5a2a4f15e9c6",
        "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
        "task_type": "query",
        "is_successful": false,
        "duration": 3025,
        "output": "eyJJRHMiOlsxLDIsMyw0LDUsNiw3LDgsOSwxMCwxMSwxMl0sInRvdGFsQW1vdW50Ijo2NzYxLjQ1MDAwMDAwMDAwMSwidG90YWxBbW91bnRJbkVVUiI6MjUxMC4xOH0=",
        "metric_value": 0,
        "executed_at": "2023-07-28T14:45:30Z"
    }
]};*/
    }

    getSlowTaskRuns(uuid, dateFrom, dateTo, threshold, taskType) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL 
                + 'organization/' + uuid 
                + '/tasks/slow/' + dateFrom 
                + '/' + dateTo 
                + '/' + threshold 
                + '/' + taskType,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {"data": [
            {
                "pipeline_run_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
                "task_uuid": "31647efa-5084-4b5a-ae06-0da647389e02",
                "task_type": "aggregator",
                "pipeline_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
                "pipeline_type": "generic",
                "is_successful": true,
                "duration": 0,
                "output": "",
                "metric_value": 13500,
                "executed_at": "2024-01-13T18:36:58Z"
            },
            {
                "pipeline_run_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
                "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
                "task_type": "query",
                "pipeline_uuid": "9c206d84-da49-49a1-a7f0-05e4f469f90d",
                "pipeline_type": "metric",
                "is_successful": true,
                "duration": 639,
                "output": "",
                "metric_value": 0,
                "executed_at": "2024-01-13T18:36:57Z"
            },
            {
                "pipeline_run_uuid": "e2e0a319-c118-44c4-bcec-e3f63a0d5d6c",
                "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
                "task_type": "query",
                "is_successful": false,
                "pipeline_uuid": "e2e0a319-c118-44c4-bcec-e3f63a0d5d6c",
                "duration": 3024,
                "output": "",
                "metric_value": 0,
                "executed_at": "2023-07-28T14:45:40Z"
            },
            {
                "pipeline_run_uuid": "a1cbc9b9-fc34-4e56-91c4-5a2a4f15e9c6",
                "task_uuid": "15fadda7-7a14-49b0-9eb2-239f21eece93",
                "task_type": "query",
                "is_successful": false,
                "pipeline_uuid": "e2e0a319-c118-44c4-bcec-e3f63a0d5d6c",
                "duration": 3025,
                "output": "",
                "metric_value": 0,
                "executed_at": "2023-07-28T14:45:30Z"
            }
        ]};*/
    }

    getPipelineValues(uuid, dateFrom, dateTo) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/values/' + dateFrom + '/' + dateTo,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {
            data: {
                "num_of_successes": 3567,
                "num_of_failures": 98,
                "average_duration": 5678,
                "is_last_run_successful": false,
                "last_run_duration": 12345,
                "last_run_executed_at": "2021-01-05 05:27:00"
            }
        };*/
    }

    getLastPipelineValues(uuid, num = 5) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/last-values/' + num,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {
            data: [
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": "jjj",
                    "executed_at": "2023-07-28 18:45:54"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1400,
                    "executed_at": "2023-07-28 17:15:01"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1461.454,
                    "executed_at": "2023-07-28 16:15:02"
                },
                {
                    "pipeline_uuid": "d63f9eca-29e9-4d99-a1d0-7f2eb0be0a38",
                    "metric_value": 1200.01,
                    "executed_at": "2023-07-28 15:15:02"
                }
            ]
        };*/
    }

    getTaskStats(uuid, dateFrom, dateTo) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'tasks/' + uuid + '/stats/' + dateFrom + '/' + dateTo,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {
            data: {
                "num_of_successes": 24,
                "num_of_failures": 130,
                "average_duration": 234,
                "is_last_run_successful": true,
                "last_run_duration": 98451234,
                "last_run_executed_at": "2021-06-05 05:20:00"
            }
        };*/
    }

    getPipelineAggregatedStats(uuid, dateFrom, period, periodLength) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/aggregated-stats/' + dateFrom + '/' + period + '/' + periodLength,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );

        /*return {
            data: [
            {
                "date_from": "2021-01-01 00:00:00",
                "date_to": "2021-01-31 00:00:00",
                "num_of_successes": 25,
                "num_of_failures": 78,
                "average_duration": 234,
            },
            {
                "date_from": "2021-02-01 00:00:00",
                "date_to": "2021-02-28 00:00:00",
                "num_of_successes": 2,
                "num_of_failures": 90,
                "average_duration": 23499,
            },
            {
                "date_from": "2021-03-01 00:00:00",
                "date_to": "2021-03-31 00:00:00",
                "num_of_successes": 56,
                "num_of_failures": 13,
                "average_duration": 1454,
            },
            {
                "date_from": "2021-04-01 00:00:00",
                "date_to": "2021-04-30 00:00:00",
                "num_of_successes": 125,
                "num_of_failures": 78,
                "average_duration": 2534,
            },
            {
                "date_from": "2021-05-01 00:00:00",
                "date_to": "2021-05-31 00:00:00",
                "num_of_successes": 250,
                "num_of_failures": 135,
                "average_duration": 2304,
            },
            {
                "date_from": "2021-06-01 00:00:00",
                "date_to": "2021-06-30 00:00:00",
                "num_of_successes": 98,
                "num_of_failures": 45,
                "average_duration": 903,
            },
            {
                "date_from": "2021-07-01 00:00:00",
                "date_to": "2021-07-31 00:00:00",
                "num_of_successes": 103,
                "num_of_failures": 135,
                "average_duration": 1903,
            },
            {
                "date_from": "2021-08-01 00:00:00",
                "date_to": "2021-08-30 00:00:00",
                "num_of_successes": 245,
                "num_of_failures": 748,
                "average_duration": 972,
            },
            {
                "date_from": "2021-09-01 00:00:00",
                "date_to": "2021-09-31 00:00:00",
                "num_of_successes": 225,
                "num_of_failures": 78,
                "average_duration": 5234,
            },
            {
                "date_from": "2021-10-01 00:00:00",
                "date_to": "2021-10-30 00:00:00",
                "num_of_successes": 125,
                "num_of_failures": 48,
                "average_duration": 1234,
            },
            {
                "date_from": "2021-11-01 00:00:00",
                "date_to": "2021-11-31 00:00:00",
                "num_of_successes": 425,
                "num_of_failures": 38,
                "average_duration": 5234,
            },
            {
                "date_from": "2021-12-01 00:00:00",
                "date_to": "2021-12-31 00:00:00",
                "num_of_successes": 65,
                "num_of_failures": 46,
                "average_duration": 1234,
            },
        ]};*/
    }

    getTaskAggregatedStats(uuid, dateFrom, period, periodLength) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'tasks/' + uuid + '/aggregated-stats/' + dateFrom + '/' + period + '/' + periodLength,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
        /*return {
            data: [
            {
                "date_from": "2021-01-01 00:00:00",
                "date_to": "2021-01-31 00:00:00",
                "num_of_successes": 25,
                "num_of_failures": 78,
                "average_duration": 234,
            },
            {
                "date_from": "2021-02-01 00:00:00",
                "date_to": "2021-02-28 00:00:00",
                "num_of_successes": 2,
                "num_of_failures": 90,
                "average_duration": 23499,
            },
            {
                "date_from": "2021-03-01 00:00:00",
                "date_to": "2021-03-31 00:00:00",
                "num_of_successes": 56,
                "num_of_failures": 13,
                "average_duration": 1454,
            },
            {
                "date_from": "2021-04-01 00:00:00",
                "date_to": "2021-04-30 00:00:00",
                "num_of_successes": 125,
                "num_of_failures": 78,
                "average_duration": 2534,
            },
            {
                "date_from": "2021-05-01 00:00:00",
                "date_to": "2021-05-31 00:00:00",
                "num_of_successes": 250,
                "num_of_failures": 135,
                "average_duration": 2304,
            },
            {
                "date_from": "2021-06-01 00:00:00",
                "date_to": "2021-06-30 00:00:00",
                "num_of_successes": 98,
                "num_of_failures": 45,
                "average_duration": 903,
            },
            {
                "date_from": "2021-07-01 00:00:00",
                "date_to": "2021-07-31 00:00:00",
                "num_of_successes": 103,
                "num_of_failures": 135,
                "average_duration": 1903,
            },
            {
                "date_from": "2021-08-01 00:00:00",
                "date_to": "2021-08-30 00:00:00",
                "num_of_successes": 245,
                "num_of_failures": 748,
                "average_duration": 972,
            },
            {
                "date_from": "2021-09-01 00:00:00",
                "date_to": "2021-09-31 00:00:00",
                "num_of_successes": 225,
                "num_of_failures": 78,
                "average_duration": 5234,
            },
            {
                "date_from": "2021-10-01 00:00:00",
                "date_to": "2021-10-30 00:00:00",
                "num_of_successes": 125,
                "num_of_failures": 48,
                "average_duration": 1234,
            },
            {
                "date_from": "2021-11-01 00:00:00",
                "date_to": "2021-11-31 00:00:00",
                "num_of_successes": 425,
                "num_of_failures": 38,
                "average_duration": 5234,
            },
            {
                "date_from": "2021-12-01 00:00:00",
                "date_to": "2021-12-31 00:00:00",
                "num_of_successes": 65,
                "num_of_failures": 46,
                "average_duration": 1234,
            },
        ]};*/
    }

    getPipelineLastRun(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'pipelines/' + uuid + '/last-run/stats',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getTaskLastRun(uuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                STATS_API_URL + 'tasks/' + uuid + '/last-run/stats',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }
}

export default new StatService();
