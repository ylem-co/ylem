import useBreadcrumbs from "use-react-router-breadcrumbs";
import { Link } from "react-router-dom";
import routes from "../routes";

const Breadcrumbs = () => {

  const breadcrumbs = useBreadcrumbs(routes, { disableDefaults: true });

  return (
    <>
      {breadcrumbs.map(({ match, breadcrumb }, index) => (
        <span className="bc" key={match.url}>
          <Link to={match.url || ""}>{breadcrumb}</Link>
          {index < breadcrumbs.length - 1 && " > "}
        </span>
      ))}
    </>
  );
};

export default Breadcrumbs;
