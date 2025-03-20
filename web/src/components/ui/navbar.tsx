import { useRef, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { GitHubLogo } from "@/components/logos/github";
import { Package, Terminal } from "lucide-react";
import { motion } from "motion/react";
import { Link } from "@tanstack/react-router";

export function Navbar() {
	const navRef = useRef<HTMLElement>(null);
	const [scrolled, setScrolled] = useState(false);

	useEffect(() => {
		const handleScroll = () => {
			setScrolled((scrolled) => {
				const isScrolled = window.scrollY > 10;
				if (isScrolled !== scrolled) {
					return isScrolled;
				}
				return scrolled;
			});
		};

		window.addEventListener("scroll", handleScroll);
		return () => window.removeEventListener("scroll", handleScroll);
	}, []);

	return (
		<motion.header
			ref={navRef}
			animate={scrolled ? "scrolled" : "initial"}
			initial={{
				paddingInline: "calc(var(--spacing) * 4)",
				paddingBlock: "calc(var(--spacing) * 6)",
				boxShadow:
					"var(--tw-inset-shadow), var(--tw-inset-ring-shadow), var(--tw-ring-offset-shadow), var(--tw-ring-shadow), var(--tw-shadow)",
			}}
			className="fixed top-0 left-0 right-0 z-50 flex bg-background"
			variants={{
				initial: {
					paddingInline: "calc(var(--spacing) * 4)",
					paddingBlock: "calc(var(--spacing) * 6)",
					boxShadow:
						"var(--tw-inset-shadow), var(--tw-inset-ring-shadow), var(--tw-ring-offset-shadow), var(--tw-ring-shadow), var(--tw-shadow)",
				},
				scrolled: {
					paddingInline: "calc(var(--spacing) * 4)",
					paddingBlock: "calc(var(--spacing) * 4)",
					boxShadow:
						"var(--tw-inset-shadow), var(--tw-inset-ring-shadow), var(--tw-ring-offset-shadow), var(--tw-ring-shadow), var(--tw-shadow)",
					"--tw-shadow":
						"0 1px 3px 0 var(--tw-shadow-color, rgb(0 0 0 / 0.1)), 0 1px 2px -1px var(--tw-shadow-color, rgb(0 0 0 / 0.1))",
				},
			}}
			transition={{ type: "spring", stiffness: 300, damping: 30 }}
		>
			<div className="container mx-auto px-4 flex justify-between items-center">
				<div className="flex items-center gap-2">
					<Link
						to="/"
						className="text-2xl font-bold flex items-center gap-2 transition-colors hover:text-primary"
					>
						<Package className="h-6 w-6" /> Konbini
					</Link>
				</div>
				<div className="flex items-center gap-4">
					<nav className="flex gap-6">
						<Button variant={scrolled ? "default" : "outline"} asChild>
							<Link to="/login">
								<Terminal />
								Login
							</Link>
						</Button>
					</nav>
					<div className="flex items-center gap-4">
						<Button asChild variant="outline" size="icon">
							<a
								href="https://github.com/juancwu/konbini"
								target="_blank"
								rel="noreferrer"
							>
								<GitHubLogo className="w-4 h-4" />
								<span className="sr-only">GitHub</span>
							</a>
						</Button>
					</div>
				</div>
			</div>
		</motion.header>
	);
}
