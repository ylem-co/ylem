import React from 'react'
import {Link} from 'react-router-dom';

const Page404 = () => {
  return (
    <>
      <h1 className="float-left display-3 mr-4">404</h1>
      <h4 className="pt-3">The page you were looking for was not found.</h4>
      <p className="text-muted float-left">
        But don't feel lost and forgotten, we are always there for you. Please go to the&nbsp;<Link to="/">main page</Link> and start from there.
      </p>
    </>
  )
}

export default Page404
