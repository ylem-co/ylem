import React from 'react';

export const releaseNotes = [
    {
        //selector: '.tour-step-profiling',
        content: () => (
            <div>
                <h2 className="card-title mb-4 alternativeHeader">Pipelines instead of pipelines!</h2> 
                <div>
                    Evolution is a continuous process. If we look at the functionality of our pipelines when Ylem just started and compare it with the status quo the difference is like the distance from Earth to the Moon.
                    <br/><br/>
                    Does not it just represent simple business processes anymore, but gives you a Swiss Army knife for data streaming, transformation, enriching, ingestion, cleaning, orchestration, and so on and so on.
                    <br/><br/>
                    Therefore, we took one more step further and decided to rename pipelines to pipelines, which will allow us to introduce even more exciting features shortly.
                    <br/><br/>
                </div>
                <div className="text-center px-4">
                    <img alt="Pipelines instead of pipelines" src="/images/release-notes/r17.png" width="600px"/>
                </div>
            </div>
        ),
        //position: [50, 30],
    },
]
