import { useRef, useEffect, useState } from "react";
import { ModeToggle } from "@/components/ui/mode-toggle";
import { Button } from "@/components/ui/button";
import { GitHubLogo } from "@/components/logos/github";
import { Package } from "lucide-react";
import { motion } from "motion/react";

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
			className="fixed top-0 left-0 right-0 z-50 flex bg-background"
			variants={{
				initial: {
					paddingInline: "calc(var(--spacing) * 4)",
					paddingBlock: "calc(var(--spacing) * 6)",
					boxShadow: "none",
				},
				scrolled: {
					paddingInline: "calc(var(--spacing) * 4)",
					paddingBlock: "calc(var(--spacing) * 4)",
					boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)",
				},
			}}
			transition={{ type: "spring", stiffness: 300, damping: 30, mass: 0.8 }}
		>
			<div className="container mx-auto px-4 flex justify-between items-center">
				<div className="flex items-center gap-2">
					<span className="text-2xl font-bold flex items-center gap-2">
						<Package className="h-6 w-6" /> Konbini
					</span>
				</div>
				<div className="flex items-center gap-4">
					<nav className="hidden md:flex gap-6">
						<a
							href="#features"
							className="hover:text-primary transition-colors"
						>
							Features
						</a>
						<a
							href="#what-is-bento"
							className="hover:text-primary transition-colors"
						>
							What is a Bento?
						</a>
						<a
							href="#security"
							className="hover:text-primary transition-colors"
						>
							Security
						</a>
						<a
							href="#getting-started"
							className="hover:text-primary transition-colors"
						>
							Getting Started
						</a>
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
						<ModeToggle />
					</div>
				</div>
			</div>
		</motion.header>
	);
}
