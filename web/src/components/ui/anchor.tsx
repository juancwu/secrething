import { Anchor as MantineAnchor } from "@mantine/core";
import { Link } from "@tanstack/react-router";

import type { LinkProps } from "@tanstack/react-router";

export function Anchor(
	props: Pick<LinkProps, "to" | "children" | "replace" | "from" | "search">,
) {
	return (
		//@ts-expect-error
		<MantineAnchor component={Link} {...props}>
			{props.children}
		</MantineAnchor>
	);
}
