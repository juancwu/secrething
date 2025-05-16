import { Anchor as MantineAnchor } from "@mantine/core";
import { Link } from "@tanstack/react-router";

import type { LinkProps } from "@tanstack/react-router";

export function Anchor(
	props: Pick<LinkProps, "to" | "children" | "replace" | "from">,
) {
	return (
		<MantineAnchor
			component={Link}
			to={props.to}
			replace={props.replace}
			from={props.from}
		>
			{props.children}
		</MantineAnchor>
	);
}
