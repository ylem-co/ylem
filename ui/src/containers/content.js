import React, { Suspense } from 'react'
import {
  Route,
  Routes
} from 'react-router-dom'

// routes config
import routes from '../routes'
  
const loading = (
  <div className="pt-3 text-center">
    <div className="sk-spinner sk-spinner-pulse"></div>
  </div>
)

const Content = () => {
  return (
    <main className="c-main content">
      <div>
        <Suspense fallback={loading}>
          <Routes>
            {routes.map((route, idx) => {
              return route.element && (
                <Route
                  key={idx}
                  path={route.path}
                  name={route.name}
                  element={route.element}
                  type={route.type}
                />
              )
            })}
          </Routes>
        </Suspense>
      </div>
    </main>
  )
}

export default React.memo(Content)
