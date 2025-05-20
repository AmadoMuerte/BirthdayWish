import { useState } from 'react';
import {
    IconSettings,
    IconFriends,
    IconHome2,
    IconLogout,
    IconPlus,
} from '@tabler/icons-react';
import { Center, Stack, Tooltip, UnstyledButton } from '@mantine/core';
import classes from './nav_bar.module.css';
import { Link } from '@tanstack/react-router';
import logo from '../../assets/react.svg';

interface NavbarLinkProps {
    icon: typeof IconHome2;
    label: string;
    active?: boolean;
    onClick?: () => void;
}

function NavbarLink({ icon: Icon, label, active, onClick }: NavbarLinkProps) {
    return (
        <Tooltip label={label} position="right" transitionProps={{ duration: 0 }}>
            <UnstyledButton onClick={onClick} className={classes.link} data-active={active || undefined}>
                <Icon size={24} stroke={1.5} />
            </UnstyledButton>
        </Tooltip>
    );
}

const mockdata = [
    { icon: IconHome2, label: 'Wishlist', path: '/app/wishlist' },
    { icon: IconPlus, label: 'Add to wishlist', path: '/app/create' },
    { icon: IconFriends, label: 'Friends', path: '/app/friends' },
];

export function Navbar() {
    const [active, setActive] = useState(0);

    const links = mockdata.map((link, index) => (
        <Link to={link.path}>
            <NavbarLink
                {...link}
                key={link.label}
                active={index === active}
                onClick={() => setActive(index)}
            />
        </Link>
    ));

    return (
        <nav className={classes.navbar}>
            <Center>
                <img src={logo} alt="logo" />
            </Center>

            <div className={classes.navbarMain}>
                <Stack justify="center" gap={0}>
                    {links}
                </Stack>
            </div>

            <Stack justify="center" gap={0}>
                <Link to='/app/settings'>
                    <NavbarLink icon={IconSettings} label="Settings" />
                </Link>
                <a href='/logout'>
                    <NavbarLink icon={IconLogout} label="Logout" />
                </a>
            </Stack>
        </nav>
    );
}