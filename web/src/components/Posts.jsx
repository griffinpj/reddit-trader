import React from 'react';
import { Breadcrumb, Layout, Menu, theme } from 'antd';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
const { Content, Sider } = Layout;


// TODO these should be the subreddits
// For each subreddit we should have a table of posts we can show
const items2 = ['avexchange', 'photomarket', 'appleswap'].map((subreddit, index) => {
    return {
        key: subreddit,
        label: 'r/' + subreddit,
        children: [
            {
                key: `${subreddit}_buy`,
                label: 'Buying',
            },
            {
                key: `${subreddit}_sell`,
                label: 'Selling',
                // TODO only show the keys that have data
                // disabled: true
            }
        ]
    };
});

// TODO preselect the first key

function Posts() {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();
    return (
        <>
            <Breadcrumb
                items={[{ title: 'Posts' }, { title: 'Subreddits' }]}
            />
            <Layout
                style={{ padding: '24px 0', background: colorBgContainer, borderRadius: borderRadiusLG }}
            >
                <Sider style={{ background: colorBgContainer }} width={200}>
                    <Menu
                        mode="inline"
                        defaultSelectedKeys={['avexchange_sell']}
                        defaultOpenKeys={['avexchange']}
                        style={{ height: '100%' }}
                        items={items2}
                    />
                </Sider>
                <Content style={{ padding: '0 24px', minHeight: 280 }}>Content</Content>
            </Layout>
        </>
    );
}


export default Posts;
